package dlm

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/namespace"
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

func Key(key string) string {
	return namespace.Key(fmt.Sprintf("dlm/%s", key))
}
