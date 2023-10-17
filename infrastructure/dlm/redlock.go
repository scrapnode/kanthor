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

	return func(key string, expiry time.Duration) DLM {
		key = fmt.Sprintf("readlock/%s", key)
		return &redlock{key: key, mu: rs.NewMutex(key, redsync.WithExpiry(expiry))}
	}, nil
}

type redlock struct {
	key string
	mu  *redsync.Mutex
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
		return fmt.Errorf("dlocker.unlock: unable to unlock because of quorum issue | key:%s", locker.key)
	}
	return nil
}
