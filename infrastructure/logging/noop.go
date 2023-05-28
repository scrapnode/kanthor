package logging

import "go.uber.org/zap"

type noop struct {
	*zap.SugaredLogger
}

// With returns a new no-op logger.
func (logger *noop) With(args ...interface{}) Logger {
	return logger.With(args)
}

func NewNoop() Logger {
	return &noop{
		SugaredLogger: zap.NewNop().Sugar(),
	}
}
