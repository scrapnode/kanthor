package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewRedis(conf *Config, logger logging.Logger) Cache {
	logger = logger.With("cache", "redis")
	return &redis{conf: conf, logger: logger}
}

type redis struct {
	conf   *Config
	logger logging.Logger

	client *goredis.Client

	mu     sync.Mutex
	status int
}

func (cache *redis) Readiness() error {
	if cache.status == patterns.StatusDisconnected {
		return nil
	}
	if cache.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return cache.client.Ping(ctx).Err()
}

func (cache *redis) Liveness() error {
	if cache.status == patterns.StatusDisconnected {
		return nil
	}
	if cache.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return cache.client.Ping(ctx).Err()
}

func (cache *redis) Connect(ctx context.Context) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	client, err := cache.connect()
	if err != nil {
		return err
	}
	cache.client = client

	cache.status = patterns.StatusConnected
	cache.logger.Info("connected")
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

	if cache.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	cache.status = patterns.StatusDisconnected
	cache.logger.Info("disconnected")

	var returning error
	if err := cache.client.Close(); err != nil {
		returning = errors.Join(returning, err)
	}
	cache.client = nil

	return returning
}

func (cache *redis) Get(ctx context.Context, key string) ([]byte, error) {
	entry, err := cache.client.Get(ctx, Key(key)).Bytes()
	// convert error type to detect later
	if errors.Is(err, goredis.Nil) {
		return nil, ErrEntryNotFound
	}

	return entry, err
}

func (cache *redis) Set(ctx context.Context, key string, entry []byte, ttl time.Duration) error {
	return cache.client.Set(ctx, Key(key), entry, ttl).Err()
}

func (cache *redis) StringGet(ctx context.Context, key string) (string, error) {
	bytes, err := cache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (cache *redis) StringSet(ctx context.Context, key string, entry string, ttl time.Duration) error {
	return cache.Set(ctx, key, []byte(entry), ttl)
}

func (cache *redis) Exist(ctx context.Context, key string) bool {
	entry, err := cache.client.Exists(ctx, Key(key)).Result()
	return err == nil && entry > 0
}

func (cache *redis) Del(ctx context.Context, key string) error {
	return cache.client.Del(ctx, Key(key)).Err()
}

func (cache *redis) ExpireAt(ctx context.Context, key string, at time.Time) (bool, error) {
	return cache.client.ExpireAt(ctx, Key(key), at).Result()
}
