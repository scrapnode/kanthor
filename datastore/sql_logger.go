package datastore

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/logging"
	"gorm.io/gorm/logger"
)

func NewSqlLogger(log logging.Logger) logger.Interface {
	return &SqlLogger{log: log}
}

type SqlLogger struct {
	log logging.Logger
}

func (logger SqlLogger) LogMode(logger.LogLevel) logger.Interface {
	return logger
}

func (logger SqlLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Infow(msg, args...)
}
func (logger SqlLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Warnw(msg, args...)
}

func (logger SqlLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Errorw(msg, args...)
}

func (logger SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)

	sql, rows := fc()
	args := []interface{}{
		"rows", rows,
		"time", float64(elapsed.Nanoseconds()) / 1e6,
	}
	if err != nil {
		args = append(args, "error", err.Error())
	}

	logger.log.Debugw(sql, args...)
}
