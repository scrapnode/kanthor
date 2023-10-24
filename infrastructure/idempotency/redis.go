package idempotency

import (
	"context"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewRedis(conf *Config, logger logging.Logger) Idempotency {
	logger = logger.With("idempotency", "redis")
	return &redis{conf: conf, logger: logger}
}

type redis struct {
	conf   *Config
	logger logging.Logger

	client *goredis.Client

	mu     sync.Mutex
	status int
}

func (idempotency *redis) Readiness() error {
	if idempotency.status == patterns.StatusDisconnected {
		return nil
	}
	if idempotency.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return idempotency.client.Ping(ctx).Err()
}

func (idempotency *redis) Liveness() error {
	if idempotency.status == patterns.StatusDisconnected {
		return nil
	}
	if idempotency.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return idempotency.client.Ping(ctx).Err()
}

func (idempotency *redis) Connect(ctx context.Context) error {
	idempotency.mu.Lock()
	defer idempotency.mu.Unlock()

	if idempotency.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	client, err := idempotency.connect()
	if err != nil {
		return err
	}
	idempotency.client = client

	idempotency.status = patterns.StatusConnected
	idempotency.logger.Info("connected")
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

	if idempotency.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	idempotency.status = patterns.StatusDisconnected
	idempotency.logger.Info("disconnected")

	if idempotency.client != nil {
		if err := idempotency.client.Close(); err != nil {
			return err
		}
	}
	idempotency.client = nil

	return nil
}

func (idempotency *redis) Validate(ctx context.Context, key string) (bool, error) {
	key = Key(key)

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
