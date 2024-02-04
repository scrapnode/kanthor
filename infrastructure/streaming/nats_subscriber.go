package streaming

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	natscore "github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/project"
	"github.com/scrapnode/kanthor/telemetry"
	"github.com/sourcegraph/conc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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

	propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
	spanName := project.Topic("streaming.publisher.sub", subscriber.name)
	go func() {
		for {
			if !subscriber.subscription.IsValid() {
				if subscriber.status == patterns.StatusConnected {
					subscriber.logger.Error(ErrSubTerminiated.Error())
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
					subscriber.logger.Errorw(err.Error(), "wait_time", fmt.Sprintf("%dms", consumer.Config.AckWait))
				}
				continue
			}

			spanner := &telemetry.Spanner{
				Tracer:   ctx.Value(telemetry.CtxTracer).(trace.Tracer),
				Contexts: make(map[string]context.Context),
			}

			events := make(map[string]*Event)
			for _, msg := range messages {
				// MetaId is event.Id
				eventId := msg.Header.Get(MetaId)

				// telemetry tracing
				var carrier propagation.MapCarrier
				if err := json.Unmarshal([]byte(msg.Header.Get(HeaderTelemetryTrace)), &carrier); err == nil {
					spanner.Contexts[eventId] = propgator.Extract(context.Background(), carrier)

					spanner.StartWithRefId(
						spanName, eventId,
						attribute.String("streaming.publisher.engine", "nats"),
						attribute.String("event.id", eventId),
					)
				}

				event := NatsMsgToEvent(msg)
				if err := event.Validate(); err != nil {
					subscriber.logger.Errorw(err.Error(), "nats_msg", utils.Stringify(msg))
					continue
				}

				// event transformation
				events[eventId] = event

			}

			spannerctx := context.WithValue(ctx, telemetry.CtxSpanner, spanner)
			errors := handler(spannerctx, events)

			var wg conc.WaitGroup
			for _, msg := range messages {
				// MetaId is event.Id
				eventId := msg.Header.Get(MetaId)

				event := events[eventId]
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

				spanner.End(spanName)
			}
			wg.Wait()
		}
	}()

	subscriber.logger.Infow("subscribed",
		"max_retry", subscriber.conf.Subscriber.MaxRetry,
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
		AckWait: time.Minute,
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
