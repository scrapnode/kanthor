package streaming

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	natscore "github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/project"
	"github.com/sourcegraph/conc"
)

type NatsSubscriber struct {
	name   string
	conf   *Config
	logger logging.Logger

	nats         *nats
	subscription *natscore.Subscription

	mu     sync.Mutex
	status int
}

func (subscriber *NatsSubscriber) Name() string {
	return subscriber.name
}

func (subscriber *NatsSubscriber) Readiness() error {
	if subscriber.status == patterns.StatusDisconnected {
		return nil
	}
	if subscriber.status != patterns.StatusConnected {
		return ErrSubNotConnected
	}

	_, err := subscriber.nats.js.StreamInfo(subscriber.conf.Name)
	return err
}

func (subscriber *NatsSubscriber) Liveness() error {
	if subscriber.status == patterns.StatusDisconnected {
		return nil
	}
	if subscriber.status != patterns.StatusConnected {
		return ErrSubNotConnected
	}

	_, err := subscriber.nats.js.StreamInfo(subscriber.conf.Name)
	return err
}

func (subscriber *NatsSubscriber) Connect(ctx context.Context) error {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()

	if subscriber.status == patterns.StatusConnected {
		return ErrSubAlreadyConnected
	}

	subscriber.logger.Info("connected")

	subscriber.status = patterns.StatusConnected
	return nil
}

func (subscriber *NatsSubscriber) Disconnect(ctx context.Context) error {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()

	if subscriber.status != patterns.StatusConnected {
		return ErrSubNotConnected
	}
	subscriber.status = patterns.StatusDisconnected
	subscriber.logger.Info("disconnected")

	var retruning error
	if subscriber.subscription.IsValid() {
		if err := subscriber.subscription.Unsubscribe(); err != nil {
			retruning = errors.Join(retruning, err)
		}
	}

	return retruning
}

func (subscriber *NatsSubscriber) Sub(ctx context.Context, topic string, handler SubHandler) error {
	// @TODO: validate topic
	topic = project.Subject(topic)
	consumer, err := subscriber.consumer(ctx, subscriber.name, topic)
	if err != nil {
		return err
	}

	subscriber.logger.Infow(
		"initialized consumer",
		"consumer_name", consumer.Config.Name,
		"consumer_created_at", consumer.Created.Format(time.RFC3339),
		"subject", consumer.Config.FilterSubject,
	)

	subscriber.subscription, err = subscriber.nats.js.PullSubscribe(
		consumer.Config.FilterSubject,
		consumer.Config.Name,
		natscore.Bind(subscriber.conf.Name, consumer.Config.Name),
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			if !subscriber.subscription.IsValid() {
				if subscriber.status == patterns.StatusConnected {
					subscriber.logger.Error("subscription is no more valid")
				}
				return
			}

			messages, err := subscriber.subscription.Fetch(subscriber.conf.Subscriber.Concurrency)
			if err != nil {
				// the subscription is no longer available because we closed it programmatically
				if errors.Is(err, natscore.ErrBadSubscription) && subscriber.status == patterns.StatusDisconnected {
					return
				}

				if !errors.Is(err, natscore.ErrTimeout) {
					subscriber.logger.Errorw(err.Error(), "timeout", fmt.Sprintf("%dms", subscriber.conf.Subscriber.Timeout))
				}
				continue
			}
			subscriber.logger.Debugw("got messages", "request_count", subscriber.conf.Subscriber.Concurrency, "response_count", len(messages))

			events := map[string]*Event{}
			for _, msg := range messages {
				event := NatsMsgToEvent(msg)
				if err := event.Validate(); err != nil {
					subscriber.logger.Errorw(err.Error(), "nats_msg", utils.Stringify(msg))
					continue
				}
				// MetaId is event.Id
				events[msg.Header.Get(MetaId)] = event
			}

			errors := handler(events)

			var wg conc.WaitGroup
			for _, msg := range messages {
				event := events[msg.Header.Get(MetaId)]
				message := msg
				wg.Go(func() {
					if err, ok := errors[event.Id]; ok && err != nil {
						if err := message.Nak(); err != nil {
							// it's important to log entire event here to trace it in our log system
							subscriber.logger.Errorw(ErrSubNakFail.Error(), "nats_msg", utils.Stringify(msg))
						}
						return
					}

					if err := message.Ack(); err != nil {
						// it's important to log entire event here to trace it in our log system
						subscriber.logger.Errorw(ErrSubAckFail.Error(), "nats_msg", utils.Stringify(msg))
					}
				})

			}
			wg.Wait()
		}
	}()

	subscriber.logger.Infow("subscribed",
		"max_retry", subscriber.conf.Subscriber.MaxRetry,
		"timeout", subscriber.conf.Subscriber.Timeout,
		"concurrency", subscriber.conf.Subscriber.Concurrency,
	)
	return nil
}

func (subscriber *NatsSubscriber) consumer(ctx context.Context, name, topic string) (*natscore.ConsumerInfo, error) {
	conf := &natscore.ConsumerConfig{
		// common config
		Name:          name,
		FilterSubject: fmt.Sprintf("%s.>", topic),

		// advance config
		MaxDeliver: subscriber.conf.Subscriber.MaxRetry + 1,
		// buffer 25% of timeout to make sure we have time to do other stuffs
		AckWait: time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout*125/100),
		// if MaxRequestBatch is 1, and we are going to request 2, we will get an error
		MaxRequestBatch: subscriber.conf.Subscriber.Concurrency,
		// if MaxAckPending is 30000, and we are processing 29999 message already
		// then we are going to request 1000, we will only get 1
		MaxAckPending: subscriber.conf.Subscriber.Throughput,

		// internal config
		DeliverPolicy: natscore.DeliverAllPolicy,
		AckPolicy:     natscore.AckExplicitPolicy,
	}

	// verify persistent consumer
	conf.Durable = conf.Name
	consumer, err := subscriber.nats.js.ConsumerInfo(subscriber.conf.Name, conf.Name, natscore.Context(ctx))

	if err == nil {
		subscriber.nats.js.UpdateConsumer(subscriber.conf.Name, conf, natscore.Context(ctx))
		return consumer, nil
	}

	// not found, create a new one
	if errors.Is(err, natscore.ErrConsumerNotFound) {
		return subscriber.nats.js.AddConsumer(subscriber.conf.Name, conf, natscore.Context(ctx))
	}

	// otherwise return the error
	return nil, err
}
