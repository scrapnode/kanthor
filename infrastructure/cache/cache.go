package cache

import (
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
		logger.Warnw("unable to parse conf.Uri, use memory cache", "uri", conf.Uri, "component", "cache")
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

	return nil, fmt.Errorf("cache: unknow engine")
}

type Cache interface {
	patterns.Connectable
	Get(key string) ([]byte, error)
	Set(key string, entry []byte, ttl time.Duration) error
	Exist(key string) bool
	Del(key string) error
}

func Key(values ...string) string {
	return strings.Join(values, "/")
}
