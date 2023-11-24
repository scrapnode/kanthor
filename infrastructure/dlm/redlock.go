package dlm

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	wrapper "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredis "github.com/redis/go-redis/v9"
)

func NewRedlock(conf *Config) (Factory, error) {
	opts, err := goredis.ParseURL(conf.Uri)
	if err != nil {
		return nil, err
	}

	client := goredis.NewClient(opts)
	rs := redsync.New(wrapper.NewPool(client))

	return func(key string, opts ...Option) DistributedLockManager {
		key = Key(key)

		conf := &Config{Uri: conf.Uri, TimeToLive: conf.TimeToLive}
		for _, opt := range opts {
			opt(conf)
		}

		return &redlock{
			key:  key,
			conf: conf,
			mu:   rs.NewMutex(key, redsync.WithExpiry(time.Millisecond*time.Duration(conf.TimeToLive))),
		}
	}, nil
}

type redlock struct {
	key string

	conf *Config
	mu   *redsync.Mutex
}

func (locker *redlock) Lock(ctx context.Context) error {
	return locker.mu.LockContext(ctx)
}

func (locker *redlock) Unlock(ctx context.Context) error {
	ok, err := locker.mu.UnlockContext(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("dlm.unlock: unable to unlock because of quorum issue | key:%s", locker.key)
	}

	return nil
}
