package cache

import (
	"context"
	"errors"
	goredis "github.com/redis/go-redis/v9"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func NewRedis(conf *Config, logger logging.Logger) Cache {
	logger = logger.With("cache", "memory")
	return &redis{conf: conf, logger: logger}
}

type redis struct {
	conf   *Config
	logger logging.Logger

	mu     sync.Mutex
	client *goredis.Client
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
	uri, err := url.Parse(cache.conf.Uri)
	if err != nil {
		return nil, err
	}
	conf := &goredis.Options{
		Addr: uri.Host,
		DB:   0, // use default DB
	}
	db := strings.Trim(uri.Path, "/")
	if db != "" {
		if n, err := strconv.Atoi(db); err == nil {
			conf.DB = n
		} else {
			cache.logger.Warnw("unable to parse db number, use default 0", "db", db)
		}
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

func (cache *redis) Get(key string) ([]byte, error) {
	entry, err := cache.client.Get(context.Background(), key).Bytes()
	// convert error type to detect later
	if errors.Is(err, goredis.Nil) {
		return nil, ErrEntryNotFound
	}

	return entry, nil
}

func (cache *redis) Set(key string, entry []byte, ttl time.Duration) error {
	return cache.client.Set(context.Background(), key, entry, ttl).Err()
}

func (cache *redis) Exist(key string) bool {
	entry, err := cache.client.Exists(context.Background(), key).Result()
	return err == nil && entry > 0
}

func (cache *redis) Del(key string) error {
	return cache.client.Del(context.Background(), key).Err()
}
