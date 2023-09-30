package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewRedis(conf *Config, logger logging.Logger) Cache {
	logger = logger.With("cache", "redis")
	return &redis{conf: conf, logger: logger}
}

type redis struct {
	conf   *Config
	logger logging.Logger

	mu     sync.Mutex
	client *goredis.Client
}

func (cache *redis) Readiness() error {
	if cache.client == nil {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return cache.client.Ping(ctx).Err()
}

func (cache *redis) Liveness() error {
	if cache.client == nil {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return cache.client.Ping(ctx).Err()
}

func (cache *redis) Connect(ctx context.Context) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.client != nil {
		return ErrAlreadyConnected
	}

	client, err := cache.connect()
	if err != nil {
		return err
	}

	cache.logger.Info("connected")
	cache.client = client
	return nil
}

func (cache *redis) connect() (*goredis.Client, error) {
	conf, err := goredis.ParseURL(cache.conf.Uri)
	if err != nil {
		return nil, err
	}

	return goredis.NewClient(conf), nil
}

func (cache *redis) Disconnect(ctx context.Context) error {
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

func (cache *redis) Get(ctx context.Context, key string) ([]byte, error) {
	entry, err := cache.client.Get(ctx, key).Bytes()
	// convert error type to detect later
	if errors.Is(err, goredis.Nil) {
		return nil, ErrEntryNotFound
	}

	return entry, err
}

func (cache *redis) Set(ctx context.Context, key string, entry []byte, ttl time.Duration) error {
	return cache.client.Set(ctx, key, entry, ttl).Err()
}

func (cache *redis) Exist(ctx context.Context, key string) bool {
	entry, err := cache.client.Exists(ctx, key).Result()
	return err == nil && entry > 0
}

func (cache *redis) Del(ctx context.Context, key string) error {
	return cache.client.Del(ctx, key).Err()
}

func (cache *redis) ExpireAt(ctx context.Context, key string, at time.Time) (bool, error) {
	return cache.client.ExpireAt(ctx, key, at).Result()
}
