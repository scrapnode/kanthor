//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dispatcher"
	uc "github.com/scrapnode/kanthor/usecases/dispatcher"
)

func InitializeDispatcher(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		dispatcher.New,
		infrastructure.New,
		ResolveDispatcherSubscriberConfig,
		streaming.NewSubscriber,
		InitializeDispatcherUsecase,
	)
	return nil, nil
}

func InitializeDispatcherUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure) (uc.Dispatcher, error) {
	wire.Build(
		uc.New,
		ResolveDispatcherPublisherConfig,
		streaming.NewPublisher,
		ResolveDispatcherSenderConfig,
		sender.New,
	)
	return nil, nil
}

func ResolveDispatcherSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Dispatcher.Subscriber
}

func ResolveDispatcherPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &conf.Scheduler.Publisher
}

func ResolveDispatcherSenderConfig(conf *config.Config, logger logging.Logger) *sender.Config {
	return &conf.Dispatcher.Sender
}
