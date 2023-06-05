//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers"
	"github.com/scrapnode/kanthor/dataplane/services"
	confprovider "github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/ioc"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

func InitializeServer(provider confprovider.Provider) (servers.Servers, error) {
	wire.Build(
		servers.New,
		InitializeConfig,
		GetLoggingConfig,
		ioc.InitializeLogger,
		InitializeServices,
	)
	return nil, nil
}

func InitializeConfig(provider confprovider.Provider) (*config.Config, error) {
	wire.Build(config.New)
	return nil, nil
}

func InitializeServices(conf *config.Config, logger logging.Logger) (services.Services, error) {
	wire.Build(
		services.New,
		GetStreamingPublisherConfig,
		ioc.InitializeStreamingPublisher,
		services.NewMessage,
	)
	return nil, nil
}

func InitializeLogger(provider confprovider.Provider) (logging.Logger, error) {
	wire.Build(logging.New, InitializeConfig, GetLoggingConfig)
	return nil, nil
}

func GetLoggingConfig(conf *config.Config) *logging.Config {
	return &conf.Dataplane.Logger
}
func GetStreamingPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &streaming.PublisherConfig{ConnectionConfig: conf.Streaming}
}
