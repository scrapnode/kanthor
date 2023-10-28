package streaming

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	natscore "github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/project"
)

func NewNats(conf *Config, logger logging.Logger) (Stream, error) {
	logger = logger.With("streaming", "nats")
	return &nats{conf: conf, logger: logger}, nil
}

func NewNatsConn(uri string, logger logging.Logger) (*natscore.Conn, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	opts := []natscore.Option{
		natscore.Name(hostname),
		natscore.ReconnectWait(3 * time.Second),
		natscore.Timeout(3 * time.Second),
		natscore.MaxReconnects(9),
		natscore.DisconnectErrHandler(func(c *natscore.Conn, err error) {
			if err != nil {
				logger.Errorf("STREAMING.NATS.DISCONNECTED: %v", err)
				return
			}
		}),
		natscore.ReconnectHandler(func(conn *natscore.Conn) {
			logger.Warnf("STREAMING.NATS.RECONNECT: %v", conn.ConnectedUrl())
		}),
		natscore.ErrorHandler(func(c *natscore.Conn, s *natscore.Subscription, err error) {
			if err == natscore.ErrSlowConsumer {
				count, bytes, err := s.Pending()
				logger.Errorf("STREAMING.NATS.ERROR: %v | subject:%s behind: %d msgs / %d bytes", err, s.Subject, count, bytes)
				return
			}

			logger.Errorf("STREAMING.NATS.ERROR: %v", err)
		}),
	}

	return natscore.Connect(uri, opts...)
}

type nats struct {
	conf   *Config
	logger logging.Logger

	conn *natscore.Conn
	js   natscore.JetStreamContext

	mu          sync.Mutex
	status      int
	publishers  map[string]Publisher
	subscribers map[string]Subscriber
}

func (streaming *nats) Readiness() error {
	if streaming.status == patterns.StatusDisconnected {
		return nil
	}
	if streaming.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	_, err := streaming.conn.RTT()
	return err
}

func (streaming *nats) Liveness() error {
	if streaming.status == patterns.StatusDisconnected {
		return nil
	}
	if streaming.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	_, err := streaming.conn.RTT()
	return err
}

func (streaming *nats) Connect(ctx context.Context) error {
	streaming.mu.Lock()
	defer streaming.mu.Unlock()

	if streaming.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	conn, err := NewNatsConn(streaming.conf.Uri, streaming.logger)
	if err != nil {
		return err
	}
	streaming.conn = conn

	js, err := conn.JetStream()
	if err != nil {
		return err
	}
	streaming.js = js

	stream, err := streaming.stream()
	if err != nil {
		return err
	}

	streaming.status = patterns.StatusConnected
	streaming.logger.Infow(
		"connected",
		"stream_name", stream.Config.Name,
		"stream_created_at", stream.Created.Format(time.RFC3339),
	)
	return nil
}

func (streaming *nats) stream() (*natscore.StreamInfo, error) {
	_, err := streaming.js.StreamInfo(streaming.conf.Name)
	// only accept ErrStreamNotFound
	if err != nil && !errors.Is(err, natscore.ErrStreamNotFound) {
		return nil, err
	}

	// prepare configurations
	sconf := &natscore.StreamConfig{
		// non-editable
		Name:    streaming.conf.Name,
		Storage: natscore.MemoryStorage,
		// editable
		Replicas: streaming.conf.Nats.Replicas,
		// project.Subject(">") "We accept all subjects that belong to the configured project and tier
		Subjects:   []string{project.Subject(">")},
		MaxMsgs:    streaming.conf.Nats.Limits.Msgs,
		MaxMsgSize: streaming.conf.Nats.Limits.MsgBytes,
		MaxBytes:   streaming.conf.Nats.Limits.Bytes,
		MaxAge:     time.Duration(streaming.conf.Nats.Limits.Age) * time.Second,
		// hardcode
		Retention: natscore.LimitsPolicy,
		Discard:   natscore.DiscardOld,
	}

	// not found, create a new one
	if err != nil {
		return streaming.js.AddStream(sconf)
	}

	// update new changes
	return streaming.js.UpdateStream(sconf)
}

func (streaming *nats) Disconnect(ctx context.Context) error {
	streaming.mu.Lock()
	defer streaming.mu.Unlock()

	if streaming.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	streaming.status = patterns.StatusDisconnected
	streaming.logger.Info("disconnected")

	var retruning error

	if len(streaming.publishers) > 0 {
		streaming.publishers = nil
	}

	if len(streaming.subscribers) > 0 {
		streaming.subscribers = nil
	}

	if !streaming.conn.IsDraining() {
		if err := streaming.conn.Drain(); err != nil {
			retruning = errors.Join(retruning, err)
		}
	}
	if !streaming.conn.IsClosed() {
		streaming.conn.Close()
	}
	streaming.conn = nil

	streaming.js = nil

	return retruning
}

func (streaming *nats) Publisher(name string) Publisher {
	// @TODO: validate name

	streaming.mu.Lock()
	defer streaming.mu.Unlock()

	if streaming.publishers == nil {
		streaming.publishers = map[string]Publisher{}
	}

	if pub, exist := streaming.publishers[name]; exist {
		return pub
	}

	publisher := &NatsPublisher{
		name:   name,
		conf:   streaming.conf,
		logger: streaming.logger.With("streaming.publisher", name),
		nats:   streaming,
	}
	streaming.publishers[name] = publisher

	return publisher
}

func (streaming *nats) Subscriber(name string) Subscriber {
	// @TODO: validate name

	streaming.mu.Lock()
	defer streaming.mu.Unlock()

	if streaming.subscribers == nil {
		streaming.subscribers = map[string]Subscriber{}
	}

	if sub, exist := streaming.subscribers[name]; exist {
		return sub
	}

	subscriber := &NatsSubscriber{
		name:   name,
		conf:   streaming.conf,
		logger: streaming.logger.With("streaming.subscriber", name),
		nats:   streaming,
	}
	streaming.subscribers[name] = subscriber

	return subscriber
}
