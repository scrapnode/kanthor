package streaming

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"sync"
	"time"
)

func NewNatsSubscriberPulling(conf *SubscriberConfig, logger logging.Logger) Subscriber {
	logger = logger.With("streaming.subscriber", "nats.pull")
	return &NatsSubscriberPulling{conf: conf, logger: logger}
}

type NatsSubscriberPulling struct {
	conf   *SubscriberConfig
	logger logging.Logger

	mu           sync.Mutex
	conn         *nats.Conn
	js           nats.JetStreamContext
	stream       *nats.StreamInfo
	subscription *nats.Subscription
}

func (subscriber *NatsSubscriberPulling) Connect(ctx context.Context) error {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()

	conn, err := NewNats(subscriber.conf.Connection, subscriber.logger)
	if err != nil {
		return err
	}
	subscriber.conn = conn

	js, err := conn.JetStream()
	if err != nil {
		return err
	}
	subscriber.js = js

	subscriber.stream, err = NewNatsStream(subscriber.conf.Connection, js)
	if err != nil {
		return err
	}

	subscriber.logger.Infow(
		"connected",
		"stream_name", subscriber.stream.Config.Name,
		"stream_created_at", subscriber.stream.Created.Format(time.RFC3339),
	)
	return nil
}

func (subscriber *NatsSubscriberPulling) Disconnect(ctx context.Context) error {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()

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
	subscriber.stream = nil

	subscriber.logger.Info("disconnected")
	return nil
}

func (subscriber *NatsSubscriberPulling) Sub(ctx context.Context, handler SubHandler) error {
	consumer, err := subscriber.consumer(ctx)
	if err != nil {
		return err
	}

	subscriber.logger.Infow(
		"initialized consumer",
		"consumer_name", consumer.Config.Name,
		"consumer_created_at", consumer.Created.Format(time.RFC3339),
		"consumer_temporary", consumer.Config.Durable == "",
	)

	subscriber.subscription, err = subscriber.js.PullSubscribe(
		consumer.Config.FilterSubject,
		// Not like Push-Based Model, Pull-Based Model requires consumer to be a durable one
		consumer.Config.Name,
		nats.Bind(subscriber.stream.Config.Name, consumer.Config.Name),
	)
	if err != nil {
		return err
	}

	subscriber.logger.Infow("subscribed")

	for {
		if !subscriber.subscription.IsValid() {
			subscriber.logger.Warnw("subscription is not more valid")
			return nil
		}

		msgs, err := subscriber.subscription.Fetch(3, nats.MaxWait(time.Millisecond*10000))
		if err != nil {
			subscriber.logger.Warnw("no more message or timeout", "error", err.Error())
			continue
		}
		subscriber.logger.Debugw("got messages", "count", len(msgs))

		maps := map[string]string{}
		events := []Event{}
		for _, msg := range msgs {
			event := natsMsgToEvent(msg)
			if err := event.Validate(); err != nil {
				subscriber.logger.Errorw(err.Error(), "nats_msg", utils.Stringify(msg))
				if err := msg.Nak(); err != nil {
					// it's important to log entire event here because we can trace it in log
					subscriber.logger.Errorw(ErrSubNakFail.Error(), "nats_msg", utils.Stringify(msg))
				}
				continue
			}

			maps[msg.Header.Get(MetaId)] = event.Id
			events = append(events, event)
		}

		results := handler(events)

		for _, msg := range msgs {
			eventId := maps[msg.Header.Get(MetaId)]

			if err, ok := results[eventId].(error); ok && err != nil {
				if err := msg.Nak(); err != nil {
					// it's important to log entire event here because we can trace it in log
					subscriber.logger.Errorw(ErrSubNakFail.Error(), "nats_msg", utils.Stringify(msg))
				}
				continue
			}

			if err := msg.Ack(); err != nil {
				// it's important to log entire event here because we can trace it in log
				subscriber.logger.Errorw(ErrSubAckFail.Error(), "nats_msg", utils.Stringify(msg))
			}
		}
	}
}

func (subscriber *NatsSubscriberPulling) consumer(ctx context.Context) (*nats.ConsumerInfo, error) {
	// prepare configurations
	conf := &nats.ConsumerConfig{
		// pull-specific
		// if MaxWaiting is 1, and more than one sub.Fetch actions, we will get an error
		MaxWaiting: subscriber.conf.Pull.MaxWaiting,
		// if MaxAckPending is 1, and we are processing 1 message already
		// then we are going to request 2, we will only get 1
		MaxAckPending: subscriber.conf.Pull.MaxAckPending,
		// if MaxRequestBatch is 1, and we are going to request 2, we will get an error
		MaxRequestBatch: subscriber.conf.Pull.MaxRequestBatch,
		// general
		Name:          subscriber.conf.Name,
		MaxDeliver:    subscriber.conf.MaxDeliver,
		FilterSubject: subscriber.conf.FilterSubject,
		// @TODO: consider apply RateLimit

		DeliverPolicy: nats.DeliverAllPolicy,
		AckPolicy:     nats.AckExplicitPolicy,
	}

	// verify persistent consumer
	conf.Durable = subscriber.conf.Name
	consumer, err := subscriber.js.ConsumerInfo(subscriber.stream.Config.Name, subscriber.conf.Name, nats.Context(ctx))

	// ideally we should update consumer options here,
	// but nats didn't support it yet,
	if err == nil {
		return consumer, nil
	}

	// not found, create a new one
	if errors.Is(err, nats.ErrConsumerNotFound) {
		return subscriber.js.AddConsumer(subscriber.stream.Config.Name, conf, nats.Context(ctx))
	}

	// otherwise return the error
	return nil, err
}