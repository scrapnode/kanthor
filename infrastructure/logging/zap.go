package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type z struct {
	*zap.SugaredLogger
}

// With returns a new no-op logger.
func (logger *z) With(args ...interface{}) Logger {
	return &z{logger.SugaredLogger.With(args...)}
}

func NewZap(conf *Config) (Logger, error) {
	var zapConfig zap.Config

	if conf.Debug {
		// running in development mode we will use a human-readable output
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	}

	var l zapcore.Level
	if err := l.UnmarshalText([]byte(conf.Level)); err != nil {
		// if something went wrong, set to debug to get as much information as possible
		l = zap.DebugLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(l)

	logger, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("logging.zap.config.Build(): %v", err))
	}

	if conf.With != nil {
		for key, value := range conf.With {
			logger = logger.With(zap.String(key, value))
		}
	}

	return &z{logger.Sugar()}, nil
}
