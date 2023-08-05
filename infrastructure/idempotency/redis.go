package idempotency

import (
	"context"
	goredis "github.com/redis/go-redis/v9"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"sync"
	"time"
)

func NewRedis(conf *Config, logger logging.Logger) Idempotency {
	logger = logger.With("idempotency", "redis")
	return &redis{conf: conf, logger: logger}
}

type redis struct {
	conf   *Config
	logger logging.Logger

	mu     sync.Mutex
	client *goredis.Client
}

func (idempotency *redis) Connect(ctx context.Context) error {
	idempotency.mu.Lock()
	defer idempotency.mu.Unlock()

	if idempotency.client != nil {
		return ErrAlreadyConnected
	}

	client, err := idempotency.connect()
	if err != nil {
		return err
	}

	idempotency.logger.Info("connected")
	idempotency.client = client
	return nil
}

func (idempotency *redis) connect() (*goredis.Client, error) {
	conf, err := goredis.ParseURL(idempotency.conf.Uri)
	if err != nil {
		return nil, err
	}

	return goredis.NewClient(conf), nil
}

func (idempotency *redis) Disconnect(ctx context.Context) error {
	idempotency.mu.Lock()
	defer idempotency.mu.Unlock()

	if idempotency.client == nil {
		return ErrNotConnected
	}

	if err := idempotency.client.Close(); err != nil {
		return err
	}
	idempotency.client = nil

	idempotency.logger.Info("disconnected")
	return nil
}

func (idempotency *redis) Validate(ctx context.Context, key string) (bool, error) {
	key = utils.Key(idempotency.conf.Namespace, key)

	var incr *goredis.IntCmd
	_, err := idempotency.client.Pipelined(ctx, func(pipe goredis.Pipeliner) error {
		incr = pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Millisecond*time.Duration(idempotency.conf.TimeToLive))
		return nil
	})
	if err != nil {
		return false, err
	}

	return incr.Val() <= 1, nil
}