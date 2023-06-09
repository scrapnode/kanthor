// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	config2 "github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers"
	"github.com/scrapnode/kanthor/dataplane/usecases/message"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/ioc"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/timer"
)

// Injectors from wire.go:

func InitializeServer(provider config.Provider) (servers.Servers, error) {
	configConfig, err := InitializeConfig(provider)
	if err != nil {
		return nil, err
	}
	loggingConfig := ResolveLoggingConfig(configConfig)
	logger, err := ioc.InitializeLogger(loggingConfig)
	if err != nil {
		return nil, err
	}
	service, err := InitializeMessageService(configConfig, logger)
	if err != nil {
		return nil, err
	}
	serversServers := servers.New(configConfig, logger, service)
	return serversServers, nil
}

func InitializeConfig(provider config.Provider) (*config2.Config, error) {
	configConfig, err := config2.New(provider)
	if err != nil {
		return nil, err
	}
	return configConfig, nil
}

func InitializeMessageService(conf *config2.Config, logger logging.Logger) (message.Service, error) {
	timerTimer := timer.New()
	publisherConfig := ResolveStreamingPublisherConfig(conf)
	publisher, err := ioc.InitializeStreamingPublisher(publisherConfig, logger)
	if err != nil {
		return nil, err
	}
	databaseConfig := ResolveDatabaseConfig(conf)
	repositoriesRepositories := repositories.New(databaseConfig, logger, timerTimer)
	service := message.NewService(conf, logger, timerTimer, publisher, repositoriesRepositories)
	return service, nil
}

func InitializeLogger(provider config.Provider) (logging.Logger, error) {
	configConfig, err := InitializeConfig(provider)
	if err != nil {
		return nil, err
	}
	loggingConfig := ResolveLoggingConfig(configConfig)
	logger, err := logging.New(loggingConfig)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// wire.go:

func ResolveLoggingConfig(conf *config2.Config) *logging.Config {
	return &conf.Dataplane.Logger
}

func ResolveStreamingPublisherConfig(conf *config2.Config) *streaming.PublisherConfig {
	return &streaming.PublisherConfig{ConnectionConfig: conf.Streaming}
}

func ResolveDatabaseConfig(conf *config2.Config) *database.Config {
	return &conf.Database
}
