//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dispatcher"
	dispatcheruc "github.com/scrapnode/kanthor/usecases/dispatcher"
)

func InitializeDispatcher(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		dispatcher.New,
		ResolveDispatcherSubscriberConfig,
		streaming.NewSubscriber,
		InitializeDispatcherUsecase,
	)
	return nil, nil
}

func InitializeDispatcherUsecase(conf *config.Config, logger logging.Logger) (dispatcheruc.Dispatcher, error) {
	wire.Build(
		dispatcheruc.New,
		timer.New,
		ResolveDispatcherPublisherConfig,
		streaming.NewPublisher,
		ResolveDispatcherSenderConfig,
		sender.New,
		ResolveDispatcherCacheConfig,
		cache.New,
		ResolveDispatcherCircuitBreakerConfig,
		circuitbreaker.New,
	)
	return nil, nil
}

func ResolveDispatcherPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &conf.Scheduler.Publisher
}

func ResolveDispatcherSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Dispatcher.Subscriber
}

func ResolveDispatcherCacheConfig(conf *config.Config) *cache.Config {
	return &conf.Dispatcher.Cache
}

func ResolveDispatcherCircuitBreakerConfig(conf *config.Config) *circuitbreaker.Config {
	return &conf.Dispatcher.CircuitBreaker
}

func ResolveDispatcherSenderConfig(conf *config.Config, logger logging.Logger) *sender.Config {
	return &conf.Dispatcher.Sender
}
