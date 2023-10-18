package authorizator

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type watcher struct {
	conf    *CasbinWatcherConfig
	logger  logging.Logger
	subject string
	nodeid  string

	conn         *nats.Conn
	subscription *nats.Subscription

	mu sync.Mutex
}

func (w *watcher) Readiness() error {
	_, err := w.conn.RTT()
	return err
}

func (w *watcher) Liveness() error {
	_, err := w.conn.RTT()
	return err
}

func (w *watcher) Connect(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	conn, err := streaming.NewNats(streaming.ConnectionConfig{Uri: w.conf.Uri}, w.logger)
	if err != nil {
		return err
	}
	w.conn = conn

	w.logger.Info("connected")
	return nil
}

func (w *watcher) Disconnect(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.subscription.IsValid() {
		if err := w.subscription.Unsubscribe(); err != nil {
			return err
		}
	}
	w.subscription = nil

	if !w.conn.IsDraining() {
		if err := w.conn.Drain(); err != nil {
			w.logger.Error(err)
		}
	}
	if !w.conn.IsClosed() {
		w.conn.Close()
	}
	w.conn = nil

	w.logger.Info("disconnected")
	return nil
}

func (w *watcher) Run(ctx context.Context, callback func(string)) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	subscription, err := w.conn.Subscribe(w.subject, func(msg *nats.Msg) {
		nodeid := string(msg.Data)
		// ignore published node
		if nodeid == w.nodeid {
			return
		}

		w.logger.Debugw("receive changes", "nodeid", nodeid)
		callback(nodeid)
	})
	if err != nil {
		return err
	}

	w.subscription = subscription
	return nil
}
