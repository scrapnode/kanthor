package authorizator

import (
	"context"
	gocasbin "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/utils"
	"net/url"
	"sync"
)

func NewCasbin(conf *Config, logger logging.Logger) Authorizator {
	logger = logger.With("authorizator", "casbin")
	return &casbin{conf: conf, logger: logger}
}

type casbin struct {
	conf   *Config
	logger logging.Logger

	watcher *watcher
	client  *gocasbin.Enforcer
}

func (authorizator *casbin) Connect(ctx context.Context) error {
	modelUrl, err := url.Parse(authorizator.conf.Casbin.ModelUri)
	if err != nil {
		return err
	}
	policyUrl, err := url.Parse(authorizator.conf.Casbin.PolicyUri)
	if err != nil {
		return err
	}
	adapter, err := gormadapter.NewAdapter(policyUrl.Scheme, authorizator.conf.Casbin.PolicyUri, true)
	if err != nil {
		return err
	}

	client, err := gocasbin.NewEnforcer(modelUrl.Host+modelUrl.Path, adapter)
	if err != nil {
		return err
	}
	if err := client.LoadModel(); err != nil {
		return err
	}
	if err := client.LoadPolicy(); err != nil {
		return err
	}
	authorizator.client = client

	authorizator.watcher = &watcher{
		nodeid:  utils.ID("casbin"),
		conf:    &authorizator.conf.Casbin.Watcher,
		logger:  authorizator.logger,
		subject: "kanthor.authorizator.casbin.watcher",
	}
	if err := authorizator.watcher.Connect(ctx); err != nil {
		return err
	}
	if err := authorizator.client.SetWatcher(authorizator.watcher); err != nil {
		return err
	}

	authorizator.logger.Info("connected")
	return nil
}

func (authorizator *casbin) Disconnect(ctx context.Context) error {
	if err := authorizator.watcher.Disconnect(ctx); err != nil {
		authorizator.logger.Error(err)
	}
	authorizator.watcher = nil

	authorizator.client = nil
	authorizator.logger.Info("disconnected")
	return nil
}

func (authorizator *casbin) Enforce(sub, dom, obj, act string) (bool, error) {
	return authorizator.client.Enforce(sub, dom, obj, act)
}

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
	conn, err := streaming.NewNats(&streaming.ConnectionConfig{Uri: w.conf.Uri}, w.logger)
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
