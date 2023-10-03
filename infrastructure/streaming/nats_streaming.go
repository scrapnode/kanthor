package streaming

import (
	"errors"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewNats(conf ConnectionConfig, logger logging.Logger) (*nats.Conn, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	opts := []nats.Option{
		nats.Name(hostname),
		nats.ReconnectWait(3 * time.Second),
		nats.Timeout(3 * time.Second),
		nats.MaxReconnects(9),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			if err != nil {
				logger.Errorf("STREAMING.NATS.DISCONNECTED: %v", err)
				return
			}
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			logger.Warnf("STREAMING.NATS.RECONNECT: %v", conn.ConnectedUrl())
		}),
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			if err == nats.ErrSlowConsumer {
				count, bytes, err := s.Pending()
				logger.Errorf("STREAMING.NATS.ERROR: %v | subject:%s behind: %d msgs / %d bytes", err, s.Subject, count, bytes)
				return
			}

			logger.Errorf("STREAMING.NATS.ERROR: %v", err)

		}),
	}

	return nats.Connect(conf.Uri, opts...)
}

func NewNatsStream(conf ConnectionConfig, js nats.JetStreamContext) (*nats.StreamInfo, error) {
	_, err := js.StreamInfo(conf.Stream.Name)
	// only accept ErrStreamNotFound
	if err != nil && !errors.Is(err, nats.ErrStreamNotFound) {
		return nil, err
	}

	// prepare configurations
	sconf := &nats.StreamConfig{
		// non-editable
		Name:    conf.Stream.Name,
		Storage: nats.FileStorage,
		// editable
		Replicas:   conf.Stream.Replicas,
		Subjects:   conf.Stream.Subjects,
		MaxMsgs:    conf.Stream.Limits.Msgs,
		MaxMsgSize: conf.Stream.Limits.MsgBytes,
		MaxBytes:   conf.Stream.Limits.Bytes,
		MaxAge:     time.Duration(conf.Stream.Limits.Age) * time.Second,
		// hardcode
		Retention: nats.LimitsPolicy,
		Discard:   nats.DiscardOld,
	}

	// not found, create a new one
	if err != nil {
		return js.AddStream(sconf)
	}

	// update new changes
	return js.UpdateStream(sconf)
}
