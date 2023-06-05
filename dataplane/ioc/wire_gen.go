// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	config2 "github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers"
	"github.com/scrapnode/kanthor/dataplane/services"
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/ioc"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

// Injectors from wire.go:

func InitializeServer(provider config.Provider) (servers.Servers, error) {
	configConfig, err := InitializeConfig(provider)
	if err != nil {
		return nil, err
	}
	loggingConfig := GetLoggingConfig(configConfig)
	logger, err := ioc.InitializeLogger(loggingConfig)
	if err != nil {
		return nil, err
	}
	services, err := InitializeServices(configConfig, logger)
	if err != nil {
		return nil, err
	}
	serversServers := servers.New(configConfig, logger, services)
	return serversServers, nil
}

func InitializeConfig(provider config.Provider) (*config2.Config, error) {
	configConfig, err := config2.New(provider)
	if err != nil {
		return nil, err
	}
	return configConfig, nil
}

func InitializeServices(conf *config2.Config, logger logging.Logger) (services.Services, error) {
	publisherConfig := GetStreamingPublisherConfig(conf)
	publisher, err := ioc.InitializeStreamingPublisher(publisherConfig, logger)
	if err != nil {
		return nil, err
	}
	message := services.NewMessage(conf, logger, publisher)
	servicesServices := services.New(logger, message)
	return servicesServices, nil
}

func InitializeLogger(provider config.Provider) (logging.Logger, error) {
	configConfig, err := InitializeConfig(provider)
	if err != nil {
		return nil, err
	}
	loggingConfig := GetLoggingConfig(configConfig)
	logger, err := logging.New(loggingConfig)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// wire.go:

func GetLoggingConfig(conf *config2.Config) *logging.Config {
	return &conf.Dataplane.Logger
}

func GetStreamingPublisherConfig(conf *config2.Config) *streaming.PublisherConfig {
	return &streaming.PublisherConfig{ConnectionConfig: conf.Streaming}
}
