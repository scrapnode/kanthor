package cache

import (
	"github.com/allegro/bigcache/v3"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewMemoryLogger(logger logging.Logger) bigcache.Logger {
	return &MemoryLogger{log: logger}
}

type MemoryLogger struct {
	log logging.Logger
}

func (logger *MemoryLogger) Printf(format string, v ...interface{}) {
	logger.log.Infof(format, v...)
}
