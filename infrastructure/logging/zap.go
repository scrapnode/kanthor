package logging

import (
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type z struct {
	*zap.SugaredLogger
}

// With returns a new no-op logger.
func (logger *z) With(args ...interface{}) Logger {
	return logger.With(args)
}

func NewZap(provider config.Provider) Logger {
	var zapConfig zap.Config

	cfg, err := GetConfig(provider)
	if err != nil {
		panic(fmt.Sprintf("logging.GetConfig(): %v", err))
	}

	if cfg.Debug {
		// running in development mode we will use a human-readable output
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	}

	var l zapcore.Level
	if err := l.UnmarshalText([]byte(cfg.Level)); err != nil {
		// if something went wrong, set to debug to get as much information as possible
		l = zap.DebugLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(l)

	logger, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("logging.zap.config.Build(): %v", err))
	}

	return &z{logger.Sugar()}
}
