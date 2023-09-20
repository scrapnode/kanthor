package streaming

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
)

func NewNatsSubscriberPushing(conf *SubscriberConfig, logger logging.Logger) Subscriber {
	logger = logger.With("streaming.subscriber", "nats.pushing")
	return &NatsSubscriberPushing{conf: conf, logger: logger}
}

type NatsSubscriberPushing struct {
	conf   *SubscriberConfig
	logger logging.Logger

	mu           sync.Mutex
	conn         *nats.Conn
	js           nats.JetStreamContext
	stream       *nats.StreamInfo
	subscription *nats.Subscription
}

func (subscriber *NatsSubscriberPushing) Connect(ctx context.Context) error {
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

func (subscriber *NatsSubscriberPushing) Disconnect(ctx context.Context) error {
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

	subscriber.logger.Info("disconnected")
	return nil
}

func (subscriber *NatsSubscriberPushing) Sub(ctx context.Context, handler SubHandler) error {
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

	subscriber.subscription, err = subscriber.conn.QueueSubscribe(
		subscriber.conf.Push.DeliverSubject,
		subscriber.conf.Push.DeliverGroup,
		func(msg *nats.Msg) {
			event := natsMsgToEvent(msg)
			if err := event.Validate(); err != nil {
				subscriber.logger.Errorw(err.Error(), "nats_msg", utils.Stringify(msg))
				if err := msg.Nak(); err != nil {
					// it's important to log entire event here because we can trace it in log
					subscriber.logger.Errorw(ErrSubNakFail.Error(), "nats_msg", utils.Stringify(msg))
				}
				return
			}

			errs := handler([]*Event{event})
			// if we got error from handler, we should retry it by no-ack action
			if err, ok := errs[event.Id]; ok && err != nil {
				if err := msg.Nak(); err != nil {
					// it's important to log entire event here because we can trace it in log
					subscriber.logger.Errorw(ErrSubNakFail.Error(), "nats_msg", utils.Stringify(msg))
				}
				return
			}

			if err := msg.Ack(); err != nil {
				// it's important to log entire event here because we can trace it in log
				subscriber.logger.Errorw(ErrSubAckFail.Error(), "nats_msg", utils.Stringify(msg))
			}
		},
	)
	if err != nil {
		return err
	}

	subscriber.logger.Infow("subscribed",
		"subscription_push_delivery_subject", subscriber.conf.Push.DeliverSubject,
		"subscription_push_delivery_group", subscriber.conf.Push.DeliverGroup,
	)

	c := make(chan int)
	<-c
	return nil
}

func (subscriber *NatsSubscriberPushing) consumer(ctx context.Context) (*nats.ConsumerInfo, error) {
	// prepare configurations
	conf := &nats.ConsumerConfig{
		// push-specific
		DeliverSubject: subscriber.conf.Push.DeliverSubject,
		DeliverGroup:   subscriber.conf.Push.DeliverGroup,
		// general
		Name:          subscriber.conf.Name,
		MaxDeliver:    subscriber.conf.MaxDeliver,
		FilterSubject: subscriber.conf.FilterSubject,
		// @TODO: consider apply RateLimit

		DeliverPolicy: nats.DeliverNewPolicy,
		AckPolicy:     nats.AckExplicitPolicy,
	}

	// do magic work to make create temporary consumer easier
	if subscriber.conf.Push.Temporary {
		// add temporary consumer
		return subscriber.js.AddConsumer(subscriber.stream.Config.Name, conf, nats.Context(ctx))
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
