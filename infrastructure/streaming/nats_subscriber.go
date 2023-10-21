package streaming

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/namespace"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/sourcegraph/conc"
)

func NewNatsSubscriber(conf *Config, logger logging.Logger) Subscriber {
	logger = logger.With("streaming.subscriber", "nats.pull")
	return &NatsSubscriber{stream: namespace.Name(conf.Name), conf: conf, logger: logger, status: patterns.StatusNone}
}

type NatsSubscriber struct {
	stream string
	conf   *Config
	logger logging.Logger
	status int

	mu           sync.Mutex
	conn         *nats.Conn
	js           nats.JetStreamContext
	subscription *nats.Subscription
}

func (subscriber *NatsSubscriber) Readiness() error {
	if subscriber.js == nil {
		return ErrNotConnected
	}

	_, err := subscriber.js.StreamInfo(subscriber.stream)
	return err
}

func (subscriber *NatsSubscriber) Liveness() error {
	if subscriber.js == nil {
		return ErrNotConnected
	}

	_, err := subscriber.js.StreamInfo(subscriber.stream)
	return err
}

func (subscriber *NatsSubscriber) Connect(ctx context.Context) error {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()

	conn, err := NewNats(subscriber.conf, subscriber.logger)
	if err != nil {
		return err
	}
	subscriber.conn = conn

	js, err := conn.JetStream(nats.Domain(namespace.Namespace()))
	if err != nil {
		return err
	}
	subscriber.js = js

	stream, err := NewNatsStream(subscriber.stream, &subscriber.conf.Nats, js)
	if err != nil {
		return err
	}

	subscriber.logger.Infow(
		"connected",
		"stream_name", stream,
		"stream_created_at", stream.Created.Format(time.RFC3339),
	)
	subscriber.status = patterns.StatusConnected
	return nil
}

func (subscriber *NatsSubscriber) Disconnect(ctx context.Context) error {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()
	subscriber.status = patterns.StatusDisconnected

	if subscriber.subscription.IsValid() {
		if err := subscriber.subscription.Unsubscribe(); err != nil {
			return err
		}
	}
	subscriber.subscription = nil

	if !subscriber.conn.IsDraining() {
		if err := subscriber.conn.Drain(); err != nil {
			subscriber.logger.Error(err)
		}
	}
	if !subscriber.conn.IsClosed() {
		subscriber.conn.Close()
	}
	subscriber.conn = nil

	subscriber.js = nil

	subscriber.logger.Info("disconnected")
	return nil
}

func (subscriber *NatsSubscriber) Sub(ctx context.Context, name, topic string, handler SubHandler) error {
	topic = namespace.Subject(topic)
	consumer, err := subscriber.consumer(ctx, name, topic)
	if err != nil {
		return err
	}

	subscriber.logger.Infow(
		"initialized consumer",
		"consumer_name", consumer.Config.Name,
		"consumer_created_at", consumer.Created.Format(time.RFC3339),
	)

	subscriber.subscription, err = subscriber.js.PullSubscribe(
		consumer.Config.FilterSubject,
		consumer.Config.Name,
		nats.Bind(subscriber.stream, consumer.Config.Name),
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			if !subscriber.subscription.IsValid() {
				subscriber.logger.Warnw("subscription is no more valid")
				return
			}

			messages, err := subscriber.subscription.Fetch(subscriber.conf.Subscriber.Concurrency)
			if err != nil {
				// the subscription is no longer available because we closed it programmatically
				if errors.Is(err, nats.ErrBadSubscription) && subscriber.status == patterns.StatusDisconnected {
					return
				}

				if !errors.Is(err, nats.ErrTimeout) {
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

func (subscriber *NatsSubscriber) consumer(ctx context.Context, name, topic string) (*nats.ConsumerInfo, error) {
	conf := &nats.ConsumerConfig{
		// common config
		Name:          name,
		FilterSubject: fmt.Sprintf("%s.>", topic),

		// advance config
		MaxDeliver: subscriber.conf.Subscriber.MaxRetry,
		// buffer 10% of timeout to make sure we have time to do other stuffs
		AckWait: time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout*110/100),
		// if MaxRequestBatch is 1, and we are going to request 2, we will get an error
		MaxRequestBatch: subscriber.conf.Subscriber.Concurrency,

		// internal config
		DeliverPolicy: nats.DeliverAllPolicy,
		AckPolicy:     nats.AckExplicitPolicy,
	}

	// verify persistent consumer
	conf.Durable = conf.Name
	consumer, err := subscriber.js.ConsumerInfo(subscriber.stream, conf.Name, nats.Context(ctx))

	if err == nil {
		subscriber.js.UpdateConsumer(subscriber.stream, conf, nats.Context(ctx))
		return consumer, nil
	}

	// not found, create a new one
	if errors.Is(err, nats.ErrConsumerNotFound) {
		return subscriber.js.AddConsumer(subscriber.stream, conf, nats.Context(ctx))
	}

	// otherwise return the error
	return nil, err
}
