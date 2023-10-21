package streaming

import (
	"context"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/namespace"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/sourcegraph/conc/pool"
)

func NewNatsPublisher(conf *Config, logger logging.Logger) Publisher {
	logger = logger.With("streaming.publisher", "nats")
	return &NatsPublisher{stream: namespace.Name(conf.Name), conf: conf, logger: logger}
}

type NatsPublisher struct {
	stream string
	conf   *Config
	logger logging.Logger

	mu   sync.Mutex
	conn *nats.Conn
	js   nats.JetStreamContext
}

func (publisher *NatsPublisher) Readiness() error {
	if publisher.js == nil {
		return ErrNotConnected
	}

	_, err := publisher.js.StreamInfo(publisher.stream)
	return err
}

func (publisher *NatsPublisher) Liveness() error {
	if publisher.js == nil {
		return ErrNotConnected
	}

	_, err := publisher.js.StreamInfo(publisher.stream)
	return err
}

func (publisher *NatsPublisher) Connect(ctx context.Context) error {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	conn, err := NewNats(publisher.conf, publisher.logger)
	if err != nil {
		return err
	}
	publisher.conn = conn

	js, err := conn.JetStream()
	if err != nil {
		return err
	}
	publisher.js = js

	stream, err := NewNatsStream(publisher.stream, &publisher.conf.Nats, js)
	if err != nil {
		return err
	}

	publisher.logger.Infow(
		"connected",
		"stream_name", stream.Config.Name,
		"stream_created_at", stream.Created.Format(time.RFC3339),
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

func (publisher *NatsPublisher) Pub(ctx context.Context, events map[string]*Event) map[string]error {
	returning := safe.Map[error]{}

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(publisher.conf.Publisher.Timeout))
	defer cancel()

	p := pool.New().WithMaxGoroutines(publisher.conf.Publisher.RateLimit)
	for key, event := range events {
		if err := event.Validate(); err != nil {
			returning.Set(key, err)
			continue
		}

		msg := NatsMsgFromEvent(event.Subject, event)
		p.Go(func() {
			ack, err := publisher.js.PublishMsg(msg, nats.Context(ctx), nats.MsgId(event.Id))
			if err != nil {
				returning.Set(key, err)
				return
			}

			publisher.logger.Debugw("published message", "msg_seq", ack.Sequence)
		})
	}
	p.Wait()

	return returning.Data()
}
