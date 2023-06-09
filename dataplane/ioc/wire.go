//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers"
	"github.com/scrapnode/kanthor/dataplane/usecases/message"
	"github.com/scrapnode/kanthor/domain/repositories"
	confprovider "github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/ioc"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/timer"
)

func InitializeServer(provider confprovider.Provider) (servers.Servers, error) {
	wire.Build(
		servers.New,
		InitializeConfig,
		ResolveLoggingConfig,
		ioc.InitializeLogger,
		InitializeMessageService,
	)
	return nil, nil
}

func InitializeConfig(provider confprovider.Provider) (*config.Config, error) {
	wire.Build(config.New)
	return nil, nil
}

func ResolveLoggingConfig(conf *config.Config) *logging.Config {
	return &conf.Dataplane.Logger
}

func ResolveStreamingPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &streaming.PublisherConfig{ConnectionConfig: conf.Streaming}
}

func ResolveDatabaseConfig(conf *config.Config) *database.Config {
	return &conf.Database
}

func InitializeMessageService(conf *config.Config, logger logging.Logger) (message.Service, error) {
	wire.Build(
		message.NewService,
		ResolveDatabaseConfig,
		repositories.New,
		timer.New,
		ResolveStreamingPublisherConfig,
		ioc.InitializeStreamingPublisher,
	)
	return nil, nil
}

func InitializeLogger(provider confprovider.Provider) (logging.Logger, error) {
	wire.Build(logging.New, InitializeConfig, ResolveLoggingConfig)
	return nil, nil
}
