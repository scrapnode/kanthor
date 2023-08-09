package streaming

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"sync"
	"time"
)

func NewNatsPublisher(conf *PublisherConfig, logger logging.Logger) Publisher {
	logger = logger.With("streaming.publisher", "nats")
	return &NatsPublisher{conf: conf, logger: logger}
}

type NatsPublisher struct {
	conf   *PublisherConfig
	logger logging.Logger

	mu   sync.Mutex
	conn *nats.Conn
	js   nats.JetStreamContext
}

func (publisher *NatsPublisher) Connect(ctx context.Context) error {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	conn, err := NewNats(publisher.conf.Connection, publisher.logger)
	if err != nil {
		return err
	}
	publisher.conn = conn

	js, err := conn.JetStream()
	if err != nil {
		return err
	}
	publisher.js = js

	stream, err := NewNatsStream(publisher.conf.Connection, js)
	if err != nil {
		return err
	}

	publisher.logger.Infow(
		"connected",
		"stream_name", stream.Config.Name, "stream_created_at", stream.Created.Format(time.RFC3339),
	)
	return nil
}

func (publisher *NatsPublisher) Disconnect(ctx context.Context) error {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	if !publisher.conn.IsDraining() {
		if err := publisher.conn.Drain(); err != nil {
			publisher.logger.Error(err)
		}
	}
	if !publisher.conn.IsClosed() {
		publisher.conn.Close()
	}
	publisher.conn = nil

	publisher.js = nil

	publisher.logger.Info("disconnected")
	return nil
}

func (publisher *NatsPublisher) Pub(ctx context.Context, event *Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	msg := natsMsgFromEvent(event.Subject, event)
	ack, err := publisher.js.PublishMsg(msg, nats.Context(ctx), nats.MsgId(event.Id))
	if err != nil {
		return fmt.Errorf("streaming.publisher.nats: %w", err)
	}

	publisher.logger.Debugw("published message", "msg_seq", ack.Sequence)
	return nil
}
