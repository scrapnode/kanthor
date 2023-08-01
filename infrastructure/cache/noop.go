package cache

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"time"
)

func NewNoop(conf *Config, logger logging.Logger) Cache {
	logger = logger.With("cache", "noop")
	return &noop{conf: conf, logger: logger}
}

type noop struct {
	conf   *Config
	logger logging.Logger
}

func (cache *noop) Connect(ctx context.Context) error {
	cache.logger.Info("connected")
	return nil
}

func (cache *noop) Disconnect(ctx context.Context) error {
	cache.logger.Info("disconnected")
	return nil
}

func (cache *noop) Get(key string) ([]byte, error) {
	return nil, ErrEntryNotFound
}

func (cache *noop) Set(key string, entry []byte, ttl time.Duration) error {
	return nil
}

func (cache *noop) Exist(key string) bool {
	return false
}

func (cache *noop) Del(key string) error {
	return nil
}
