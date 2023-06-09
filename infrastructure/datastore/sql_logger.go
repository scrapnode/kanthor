package datastore

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"gorm.io/gorm/logger"
	"time"
)

func NewSqlLogger(log logging.Logger) logger.Interface {
	return &GormLogger{log: log}
}

type GormLogger struct {
	log logging.Logger
}

func (logger GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return logger
}

func (logger GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Infow(msg, args...)
}
func (logger GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Warnw(msg, args...)
}

func (logger GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	logger.log.Errorw(msg, args...)
}

func (logger GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
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
