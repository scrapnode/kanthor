package cache

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewNoop(conf *Config, logger logging.Logger) Cache {
	logger = logger.With("cache", "noop")
	return &noop{conf: conf, logger: logger}
}

type noop struct {
	conf   *Config
	logger logging.Logger
}

func (cache *noop) Readiness() error {
	return nil
}

func (cache *noop) Liveness() error {
	return nil
}

func (cache *noop) Connect(ctx context.Context) error {
	cache.logger.Info("connected")
	return nil
}

func (cache *noop) Disconnect(ctx context.Context) error {
	cache.logger.Info("disconnected")
	return nil
}

func (cache *noop) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, ErrEntryNotFound
}

func (cache *noop) Set(ctx context.Context, key string, entry []byte, ttl time.Duration) error {
	return nil
}

func (cache *noop) Exist(ctx context.Context, key string) bool {
	return false
}

func (cache *noop) Del(ctx context.Context, key string) error {
	return nil
}

func (cache *noop) ExpireAt(ctx context.Context, key string, at time.Time) (bool, error) {
	return true, nil
}
