package streaming

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"strings"
	"sync"
	"time"
)

func NewNatsSubscriber(conf *SubscriberConfig, logger logging.Logger) Subscriber {
	logger = logger.With("streaming.subscriber", "nats")
	return &NatsSubscriber{conf: conf, logger: logger}
}

type NatsSubscriber struct {
	conf   *SubscriberConfig
	logger logging.Logger

	mu           sync.Mutex
	conn         *nats.Conn
	js           jetstream.JetStream
	stream       jetstream.Stream
	subscription *nats.Subscription
}

func (subscriber *NatsSubscriber) Connect(ctx context.Context) error {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()

	conn, err := NewNats(subscriber.conf.Connection, subscriber.logger)
	if err != nil {
		return err
	}
	subscriber.conn = conn

	js, err := jetstream.New(subscriber.conn)
	if err != nil {
		return err
	}
	subscriber.js = js

	subscriber.stream, err = NewNatsStream(subscriber.conf.Connection, js)
	if err != nil {
		return err
	}
	// make sure we have created it successfully
	info, err := subscriber.stream.Info(context.Background())
	if err != nil {
		return err
	}

	subscriber.logger.Infow(
		"connected",
		"stream_name", info.Config.Name, "stream_created_at", info.Created.Format(time.RFC3339),
	)
	return nil
}

func (subscriber *NatsSubscriber) Disconnect(ctx context.Context) error {
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

func (subscriber *NatsSubscriber) Sub(ctx context.Context, handler SubHandler) error {
	consumer, err := subscriber.consumer(ctx)
	if err != nil {
		return err
	}
	// make sure we have saved it successfully
	info, err := consumer.Info(ctx)
	if err != nil {
		return err
	}
	subscriber.logger.Infow(
		"initialized consumer",
		"consumer_name", info.Config.Name,
		"consumer_created_at", info.Created.Format(time.RFC3339),
		"consumer_temporary", info.Config.Durable == "",
	)

	subscriber.subscription, err = subscriber.conn.QueueSubscribe(
		subscriber.conf.Topic,
		subscriber.conf.Group,
		func(msg *nats.Msg) {
			event := subscriber.transform(msg)
			if err := event.Validate(); err != nil {
				subscriber.logger.Errorw(err.Error(), "nats_msg", utils.Stringify(msg))
				if err := msg.Nak(); err != nil {
					// it's important to log entire event here because we can trace it in log
					subscriber.logger.Errorw(ErrSubNakFail.Error(), "nats_msg", utils.Stringify(msg))
				}
				return
			}

			results := handler([]Event{event})
			// if we got error from handler, we should retry it by no-ack action
			if err, ok := results[event.Id].(error); ok && err != nil {
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

	subscriber.logger.Infow("subscribed",
		"subscription_topic", subscriber.conf.Topic,
		"subscription_group", subscriber.conf.Group,
	)
	return err
}

func (subscriber *NatsSubscriber) consumer(ctx context.Context) (jetstream.Consumer, error) {
	// prepare configurations
	conf := jetstream.ConsumerConfig{
		Name:           subscriber.conf.Name,
		DeliverSubject: subscriber.conf.Topic,
		DeliverGroup:   subscriber.conf.Group,
		MaxDeliver:     subscriber.conf.MaxDeliver,
		FilterSubject:  subscriber.conf.FilterSubject,
		// @TODO: consider apply RateLimit

		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
	}
	if conf.Name == "" {
		conf.Name = utils.MD5(subscriber.conf.Topic, subscriber.conf.Group)
	}

	// do magic work to make create temporary consumer easier
	if subscriber.conf.Temporary {
		// add temporary consumer
		return subscriber.stream.AddConsumer(ctx, conf)
	}

	// verify persistent consumer
	conf.Durable = subscriber.conf.Name
	consumer, err := subscriber.stream.Consumer(ctx, subscriber.conf.Name)

	// ideally we should update consumer options here,
	// but nats didn't support it yet,
	if err == nil {
		return consumer, nil
	}

	// not found, create a new one
	if errors.Is(err, jetstream.ErrConsumerNotFound) {
		return subscriber.stream.AddConsumer(ctx, conf)
	}

	return nil, err
}

func (subscriber *NatsSubscriber) transform(msg *nats.Msg) Event {
	event := Event{
		Subject:  msg.Subject,
		AppId:    msg.Header.Get(MetaAppId),
		Type:     msg.Header.Get(MetaType),
		Id:       msg.Header.Get(MetaId),
		Data:     msg.Data,
		Metadata: map[string]string{},
	}
	for key, value := range msg.Header {
		if strings.HasPrefix(key, "Nats") {
			continue
		}
		if key == MetaAppId || key == MetaType || key == MetaId {
			continue
		}
		event.Metadata[key] = value[0]
	}
	return event
}
