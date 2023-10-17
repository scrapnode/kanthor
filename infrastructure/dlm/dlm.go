package dlm

import (
	"context"
	"time"
)

type Factory func(key string, expiry time.Duration) DLM

type DLM interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}

func New(conf *Config) (Factory, error) {
	return NewRedlock(conf)
}
