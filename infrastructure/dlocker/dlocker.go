package dlocker

import (
	"context"
	"time"
)

type Factory func(key string, expiry time.Duration) DLocker

type DLocker interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}

var Ns = "dlocker/redlock"

func New(conf *Config) (Factory, error) {
	return NewRedlock(conf)
}
