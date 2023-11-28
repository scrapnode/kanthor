package authorizator

import (
	"context"
	"errors"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
)

type watcher struct {
	conf    *CasbinWatcherConfig
	logger  logging.Logger
	subject string
	nodeid  string

	conn         *nats.Conn
	subscription *nats.Subscription

	mu     sync.Mutex
	status int
}

func (w *watcher) Readiness() error {
	if w.status == patterns.StatusDisconnected {
		return nil
	}
	if w.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if !w.conn.IsConnected() {
		return ErrNotConnected
	}

	return nil
}

func (w *watcher) Liveness() error {
	if w.status == patterns.StatusDisconnected {
		return nil
	}
	if w.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if !w.conn.IsConnected() {
		return ErrNotConnected
	}

	return nil
}

func (w *watcher) Connect(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == patterns.StatusConnected {
		return ErrWatcherAlreadyConnected
	}

	conn, err := streaming.NewNatsConn(w.conf.Uri, w.logger)
	if err != nil {
		return err
	}
	w.conn = conn

	w.status = patterns.StatusConnected
	w.logger.Info("connected")
	return nil
}

func (w *watcher) Disconnect(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status != patterns.StatusConnected {
		return ErrWatcherNotConnected
	}
	w.status = patterns.StatusDisconnected
	w.logger.Info("disconnected")

	var returning error
	if w.subscription.IsValid() {
		if err := w.subscription.Unsubscribe(); err != nil {
			returning = errors.Join(returning, err)
		}
	}
	w.subscription = nil

	if !w.conn.IsClosed() {
		w.conn.Close()
	}
	w.conn = nil

	return returning
}

func (w *watcher) Run(ctx context.Context, callback func(string)) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	subscription, err := w.conn.Subscribe(w.subject, func(msg *nats.Msg) {
		nodeid := string(msg.Data)
		// ignore original node
		if nodeid == w.nodeid {
			w.logger.Debugw("ignore original node", "nodeid", nodeid)
			return
		}

		w.logger.Debugw("received changes", "nodeid", nodeid)
		callback(nodeid)
	})
	if err != nil {
		return err
	}
	w.logger.Debugw("running", "nodeid", w.nodeid)
	w.subscription = subscription
	return nil
}

func (w *watcher) Update() error {
	w.logger.Infow("publish changes", "subject", w.subject, "nodeid", w.nodeid)
	return w.conn.Publish(w.subject, []byte(w.nodeid))
}
