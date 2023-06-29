package streaming

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"os"
	"time"
)

func NewNats(conf *ConnectionConfig, logger logging.Logger) (*nats.Conn, error) {
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
				// @TODO: add metrics here
				logger.Error(fmt.Sprintf("got disconnected with reason: %q", err))
				return
			}
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			// @TODO: add metrics here
			logger.Infow(fmt.Sprintf("got reconnected to %v", conn.ConnectedUrl()))
		}),
	}

	return nats.Connect(conf.Uri, opts...)
}

func NewNatsStream(conf *ConnectionConfig, js jetstream.JetStream) (jetstream.Stream, error) {
	_, err := js.Stream(context.Background(), conf.Stream.Name)
	// only accept ErrStreamNotFound
	if err != nil && !errors.Is(err, jetstream.ErrStreamNotFound) {
		return nil, err
	}

	// prepare configurations
	config := jetstream.StreamConfig{
		// non-editable
		Name:    conf.Stream.Name,
		Storage: jetstream.FileStorage,
		// editable
		Replicas:   conf.Stream.Replicas,
		Subjects:   conf.Stream.Subjects,
		MaxMsgs:    conf.Stream.Limits.Msgs,
		MaxMsgSize: conf.Stream.Limits.MsgBytes,
		MaxBytes:   conf.Stream.Limits.Bytes,
		MaxAge:     time.Duration(conf.Stream.Limits.Age) * time.Second,
		// hardcode
		Retention: jetstream.LimitsPolicy,
		Discard:   jetstream.DiscardOld,
	}

	// not found, create a new one
	if err != nil {
		return js.CreateStream(context.Background(), config)
	}

	// update new changes
	return js.UpdateStream(context.Background(), config)
}
