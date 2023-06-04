package msgbroker

import (
	"context"
	"errors"
	"fmt"
	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"os"
	"strings"
	"sync"
	"time"
)

func NewNats(conf *Config, logger logging.Logger) (MsgBroker, error) {
	return &nats{conf: conf, logger: logger.With("component", "msgbroker")}, nil
}

type nats struct {
	conf   *Config
	logger logging.Logger

	mu           sync.Mutex
	client       *natsio.Conn
	js           jetstream.JetStream
	jss          jetstream.Stream
	subscription *natsio.Subscription
}

func (broker *nats) Connect(ctx context.Context) error {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	if broker.client != nil {
		return fmt.Errorf("msgbroker: %w", ErrAlreadyConnected)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	opts := []natsio.Option{
		natsio.Name(hostname),
		natsio.ReconnectWait(3 * time.Second),
		natsio.Timeout(3 * time.Second),
		natsio.MaxReconnects(9),
		natsio.DisconnectErrHandler(func(c *natsio.Conn, err error) {
			// @TODO: add metrics here
			broker.logger.Error(fmt.Sprintf("got disconnected with reason: %q", err))
		}),
		natsio.ReconnectHandler(func(conn *natsio.Conn) {
			// @TODO: add metrics here
			broker.logger.Error(fmt.Sprintf("got reconnected to %v", conn.ConnectedUrl()))
		}),
	}
	broker.client, err = natsio.Connect(broker.conf.Uri, opts...)
	if err != nil {
		return fmt.Errorf("msgbroker: %w", err)
	}

	stream, err := broker.stream(ctx)
	if err != nil {
		return fmt.Errorf("msgbroker: %w", err)
	}

	s, err := stream.Info(ctx)
	if err != nil {
		return fmt.Errorf("msgbroker: %w", err)
	}
	broker.logger.Infow("connected to a stream", "stream_name", s.Config.Name, "stream_state", s.State)

	broker.logger.Info("connected")
	return nil
}

func (broker *nats) Disconnect(ctx context.Context) error {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	if broker.client == nil {
		return fmt.Errorf("msgbroker: %w", ErrNotConnected)
	}

	if broker.subscription != nil {
		if err := broker.subscription.Drain(); err != nil {
			return fmt.Errorf("msgbroker: %w", err)
		}
	}

	if err := broker.client.Drain(); err != nil {
		return fmt.Errorf("msgbroker: %w", err)
	}

	broker.client = nil
	broker.jss = nil
	broker.subscription = nil
	broker.logger.Info("disconnected")
	return nil
}

func (broker *nats) Pub(ctx context.Context, event *Event) error {
	js, err := broker.jetstream()
	if err != nil {
		return fmt.Errorf("msgbroker: %w", err)
	}

	msg := &natsio.Msg{
		Subject: broker.subject(event),
		Header: natsio.Header{
			// nats specific
			"Nats-Msg-Id": []string{event.Id},
			// internal
			MetaTier:  []string{event.Tier},
			MetaAppId: []string{event.AppId},
			MetaType:  []string{event.Type},
			MetaId:    []string{event.Id},
		},
		Data: event.Data,
	}
	for key, value := range event.Metadata {
		msg.Header.Set(key, value)
	}
	broker.logger.Debugw("prepared message", "msg_id", event.Id, "msg_subject", msg.Subject)

	ack, err := js.PublishMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("msgbroker: %w", err)
	}

	broker.logger.Debugw("published message", "msg_sequence", ack.Sequence)
	return nil
}

func (broker *nats) Sub(ctx context.Context, handler Handler) error {
	consumer, err := broker.consumer(ctx, broker.jss)
	if err != nil {
		return fmt.Errorf("msgbroker: %w", err)
	}

	c, err := consumer.Info(ctx)
	if err != nil {
		return fmt.Errorf("msgbroker: %w", err)
	}
	broker.logger.Infow("connect to a consumer", "consumer_name", c.Name, "consumer_config", c.Config)

	subscription, err := broker.client.QueueSubscribe(c.Config.DeliverSubject, c.Config.DeliverGroup, func(msg *natsio.Msg) {
		event := &Event{
			Tier:     msg.Header.Get(MetaTier),
			AppId:    msg.Header.Get(MetaAppId),
			Type:     msg.Header.Get(MetaType),
			Id:       msg.Header.Get(MetaId),
			Data:     msg.Data,
			Metadata: map[string]string{},
		}
		for key, value := range msg.Header {
			if key == MetaTier || key == MetaAppId || key == MetaType || key == MetaId {
				continue
			}
			if strings.HasPrefix(strings.ToLower(key), "nats") {
				continue
			}

			event.Metadata[key] = value[0]
		}

		// if we got error from handler, we should retry it by no-ack action
		if err := handler(event); err != nil {
			if err := msg.Nak(); err != nil {
				// it's important to log entire event here because we can trace it in log
				broker.logger.Errorw("could not nak an event", "event", event.String())
			}
			return
		}

		if err := msg.Ack(); err != nil {
			// it's important to log entire event here because we can trace it in log
			broker.logger.Errorw("could not nak an event", "event", event.String())
		}
	})

	if err != nil {
		return fmt.Errorf("msgbroker: %w", err)
	}

	broker.subscription = subscription
	return nil
}

func (broker *nats) consumer(ctx context.Context, stream jetstream.Stream) (jetstream.Consumer, error) {
	conf := jetstream.ConsumerConfig{
		Name:           broker.conf.Consumer.Name,
		DeliverSubject: broker.conf.Consumer.Name,
		DeliverGroup:   broker.conf.Consumer.Name,
		DeliverPolicy:  jetstream.DeliverAllPolicy,
		AckPolicy:      jetstream.AckExplicitPolicy,
		// if MaxRetry is not set, we guarantee at least one MaxDeliver
		MaxDeliver: broker.conf.Consumer.MaxRetry + 1,
		// @TODO: consider apply RateLimit
	}
	if broker.conf.Consumer.FilterSubject != "" {
		conf.FilterSubject = fmt.Sprintf("%s.%s", broker.conf.Stream.Subject, broker.conf.Consumer.FilterSubject)
	}

	if !broker.conf.Consumer.Temporary {
		conf.Durable = broker.conf.Consumer.Name
		conf.InactiveThreshold = time.Hour
	}

	consumer, err := stream.Consumer(ctx, conf.Name)
	if err == nil {
		return consumer, nil
	}

	// some unexpected error was happened, report it immediately
	if err != nil && !errors.Is(err, jetstream.ErrConsumerNotFound) {
		return nil, err
	}

	return stream.AddConsumer(ctx, conf)
}

func (broker *nats) subject(event *Event) string {
	subjects := []string{
		broker.conf.Stream.Subject,
		event.Tier,
		event.AppId,
		event.Type,
	}
	return strings.Join(subjects, ".")
}

func (broker *nats) jetstream() (jetstream.JetStream, error) {
	if broker.js != nil {
		return broker.js, nil
	}

	js, err := jetstream.New(broker.client)
	if err != nil {
		return nil, err
	}

	broker.js = js
	return broker.js, nil
}

func (broker *nats) stream(ctx context.Context) (jetstream.Stream, error) {
	if broker.jss != nil {
		return broker.jss, nil
	}

	js, err := broker.jetstream()
	if err != nil {
		return nil, err
	}

	stream, err := js.Stream(ctx, broker.conf.Stream.Name)
	if errors.Is(err, jetstream.ErrStreamNotFound) {
		stream, err = js.CreateStream(ctx, jetstream.StreamConfig{
			Name:     broker.conf.Stream.Name,
			Replicas: broker.conf.Stream.Replicas,
			// only support one subject wildcard per stream
			Subjects: []string{fmt.Sprintf("%s.>", broker.conf.Stream.Subject)},
			// Retention based on the various limits that are set including: MaxMsgs, MaxBytes, MaxAge, and MaxMsgsPerSubject.
			// If any of these limits are set, whichever limit is hit first will cause the automatic deletion of the respective message(s)
			Retention: jetstream.LimitsPolicy,
			// This policy will delete the oldest messages in order to maintain the limit.
			// For example, if MaxAge is set to one minute, the server will automatically delete messages older than one minute with this policy.
			Discard:    jetstream.DiscardOld,
			MaxMsgs:    broker.conf.Stream.Limits.Msgs,
			MaxMsgSize: broker.conf.Stream.Limits.MsgBytes,
			MaxBytes:   broker.conf.Stream.Limits.Bytes,
			MaxAge:     time.Duration(broker.conf.Stream.Limits.Age) * time.Second,
		})
	}

	broker.jss = stream
	return stream, err
}
