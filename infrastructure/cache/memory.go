package cache

import (
	"context"
	"errors"
	"github.com/allegro/bigcache/v3"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"sync"
	"time"
)

func NewMemory(conf *Config, logger logging.Logger) Cache {
	logger = logger.With("component", "cache.memory")
	return &memory{conf: conf, logger: logger}
}

type memory struct {
	conf   *Config
	logger logging.Logger

	mu     sync.Mutex
	client *bigcache.BigCache
}

func (cache *memory) Connect(ctx context.Context) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.client != nil {
		return ErrAlreadyConnected
	}

	conf := bigcache.DefaultConfig(time.Duration(cache.conf.TimeToLiveInSeconds) * time.Second)
	conf.Logger = NewMemoryLogger(cache.logger)

	client, err := bigcache.New(ctx, conf)
	if err != nil {
		return err
	}

	cache.logger.Info("connected")
	cache.client = client
	return nil
}

func (cache *memory) Disconnect(ctx context.Context) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.client == nil {
		return ErrNotConnected
	}

	if err := cache.client.Close(); err != nil {
		return err
	}
	cache.client = nil

	cache.logger.Info("disconnected")
	return nil
}

func (cache *memory) Get(key string) ([]byte, error) {
	entry, err := cache.client.Get(key)
	// convert error type to detect later
	if errors.Is(err, bigcache.ErrEntryNotFound) {
		return nil, ErrEntryNotFound
	}

	return entry, nil
}

func (cache *memory) Set(key string, entry []byte, ttl time.Duration) error {
	return cache.client.Set(key, entry)
}

func (cache *memory) Exist(key string) bool {
	entry, err := cache.client.Get(key)
	return err == nil && len(entry) > 0
}

func (cache *memory) Del(key string) error {
	return cache.client.Delete(key)
}
