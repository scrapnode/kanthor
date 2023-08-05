package authorizator

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"sync"
)

type watcher struct {
	nodeid  string
	conf    *CasbinWatcherConfig
	logger  logging.Logger
	subject string

	conn         *nats.Conn
	subscription *nats.Subscription

	mu       sync.Mutex
	callback func(string)
}

func (w *watcher) Connect(ctx context.Context) error {
	conn, err := streaming.NewNats(streaming.ConnectionConfig{Uri: w.conf.Uri}, w.logger)
	if err != nil {
		return err
	}
	w.conn = conn

	subscription, err := w.conn.Subscribe(w.subject, func(msg *nats.Msg) {
		nodeid := string(msg.Data)
		// ignore publish node
		if nodeid == w.nodeid {
			return
		}

		w.logger.Debugw("receive changes", "nodeid", nodeid)
		w.callback(nodeid)
	})
	if err != nil {
		return err
	}

	w.subscription = subscription
	w.logger.Info("connected")
	return nil
}

func (w *watcher) Disconnect(ctx context.Context) error {
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

func (w *watcher) SetUpdateCallback(callback func(string)) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.callback = callback
	return nil
}

func (w *watcher) Update() error {
	return w.conn.Publish(w.subject, []byte(w.nodeid))
}

func (w *watcher) Close() {
	if err := w.Disconnect(context.Background()); err != nil {
		w.logger.Error(err)
	}
}
