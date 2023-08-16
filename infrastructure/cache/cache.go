package cache

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"net/url"
	"strings"
	"time"
)

func New(conf *Config, logger logging.Logger) (Cache, error) {
	uri, err := url.Parse(conf.Uri)
	if err != nil {
		logger.Warnw("unable to parse conf.Uri", "uri", conf.Uri)
		return nil, err
	}

	if strings.HasPrefix(uri.Scheme, "noop") {
		return NewNoop(conf, logger), nil
	}

	if strings.HasPrefix(uri.Scheme, "redis") {
		return NewRedis(conf, logger), nil
	}

	if strings.HasPrefix(uri.Scheme, "memory") {
		return NewMemory(conf, logger), nil
	}

	return nil, fmt.Errorf("cache: unknown engine")
}

type Cache interface {
	patterns.Connectable
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, entry []byte, ttl time.Duration) error
	Exist(ctx context.Context, key string) bool
	Del(ctx context.Context, key string) error
	ExpireAt(ctx context.Context, key string, at time.Time) (bool, error)
}
