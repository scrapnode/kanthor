package dlm

import (
	"context"
	"time"
)

type Factory func(key string) DistributedLockManager

type DistributedLockManager interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
	Until() time.Time
}

func New(conf *Config) (Factory, error) {
	return NewRedlock(conf)
}
