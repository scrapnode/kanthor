//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	confprovider "github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/migration/config"
	"github.com/scrapnode/kanthor/migration/operators"
)

func InitializeMigrator(provider confprovider.Provider) (operators.Operator, error) {
	wire.Build(
		operators.New,
		InitializeConfig,
		InitializeLogger,
	)
	return nil, nil
}

func InitializeConfig(provider confprovider.Provider) (*config.Config, error) {
	wire.Build(config.New)
	return nil, nil
}

func ResolveLoggingConfig(conf *config.Config) *logging.Config {
	return &conf.Migration.Logger
}

func InitializeLogger(provider confprovider.Provider) (logging.Logger, error) {
	wire.Build(logging.New, InitializeConfig, ResolveLoggingConfig)
	return nil, nil
}
