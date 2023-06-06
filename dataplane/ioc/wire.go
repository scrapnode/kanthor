//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers"
	"github.com/scrapnode/kanthor/dataplane/usecases/message"
	"github.com/scrapnode/kanthor/infrastructure/auth"
	confprovider "github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/ioc"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

func InitializeConfig(provider confprovider.Provider) (*config.Config, error) {
	wire.Build(config.New)
	return nil, nil
}

func GetLoggingConfig(conf *config.Config) *logging.Config {
	return &conf.Dataplane.Logger
}

func InitializeServer(provider confprovider.Provider) (servers.Servers, error) {
	wire.Build(
		servers.New,
		InitializeConfig,
		GetLoggingConfig,
		ioc.InitializeLogger,
		InitializeMessageUseCase,
	)
	return nil, nil
}

func GetDatastoreConfig(conf *config.Config) *datastore.Config {
	return &conf.Datastore
}

func GetAuthConfig(conf *config.Config) *auth.Config {
	return &conf.Dataplane.Auth
}

func GetStreamingPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &streaming.PublisherConfig{ConnectionConfig: conf.Streaming}
}

func InitializeMessageUseCase(conf *config.Config, logger logging.Logger) (message.Service, error) {
	wire.Build(
		message.NewService,
		GetDatastoreConfig,
		ioc.InitializeDatastore,
		message.NewRepository,
		GetAuthConfig,
		ioc.InitializeAuth,
		GetStreamingPublisherConfig,
		ioc.InitializeStreamingPublisher,
	)
	return nil, nil
}

func InitializeLogger(provider confprovider.Provider) (logging.Logger, error) {
	wire.Build(logging.New, InitializeConfig, GetLoggingConfig)
	return nil, nil
}
