package cache

import (
	"context"
	gocache "github.com/patrickmn/go-cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"sync"
	"time"
)

func NewMemory(conf *Config, logger logging.Logger) Cache {
	logger = logger.With("cache", "memory")
	return &memory{conf: conf, logger: logger}
}

type memory struct {
	conf   *Config
	logger logging.Logger

	mu     sync.Mutex
	client *gocache.Cache
}

func (cache *memory) Connect(ctx context.Context) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.client != nil {
		return ErrAlreadyConnected
	}

	ttl := time.Millisecond * time.Duration(cache.conf.TimeToLive)
	// we set cleanup interval time equal to ttl for simplify the implementation
	client := gocache.New(ttl, ttl)

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

	// delete all items to claim free memory
	cache.client.Flush()
	cache.client = nil

	cache.logger.Info("disconnected")
	return nil
}

func (cache *memory) Get(ctx context.Context, key string) ([]byte, error) {
	entry, found := cache.client.Get(key)
	if !found {
		return nil, ErrEntryNotFound
	}

	return entry.([]byte), nil
}

func (cache *memory) Set(ctx context.Context, key string, entry []byte, ttl time.Duration) error {
	cache.client.Set(key, entry, ttl)
	return nil
}

func (cache *memory) Exist(ctx context.Context, key string) bool {
	_, found := cache.client.Get(key)
	return found
}

func (cache *memory) Del(ctx context.Context, key string) error {
	cache.client.Delete(key)
	return nil
}
