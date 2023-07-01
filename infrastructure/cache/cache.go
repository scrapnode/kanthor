package cache

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"strings"
	"time"
)

func New(conf *Config, logger logging.Logger) Cache {
	return NewMemory(conf, logger)
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
