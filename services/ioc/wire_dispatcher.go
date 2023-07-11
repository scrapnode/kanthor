//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dispatcher"
	"github.com/scrapnode/kanthor/usecases"
)

func InitializeDispatcher(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		dispatcher.New,
		usecases.NewDispatcher,
		timer.New,
		ResolveDispatcherPublisherConfig,
		streaming.NewPublisher,
		ResolveDispatcherSubscriberConfig,
		streaming.NewSubscriber,
		wire.FieldsOf(new(*config.Config), "Database"),
		repositories.New,
		ResolveDispatcherSenderConfig,
		sender.New,
		ResolveDispatcherCacheConfig,
		cache.New,
		ResolveDispatcherCircuitBreakerConfig,
		circuitbreaker.New,
		ResolveDispatcherMetricConfig,
		metric.New,
	)
	return nil, nil
}

func ResolveDispatcherPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	publisher := conf.Scheduler.Publisher
	if publisher.ConnectionConfig == nil {
		publisher.ConnectionConfig = &conf.Streaming
	}
	return &publisher
}

func ResolveDispatcherSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	subscriber := conf.Dispatcher.Subscriber
	if subscriber.ConnectionConfig == nil {
		subscriber.ConnectionConfig = &conf.Streaming
	}

	return &subscriber
}

func ResolveDispatcherCacheConfig(conf *config.Config) *cache.Config {
	if conf.Dispatcher.Cache == nil {
		return &conf.Cache
	}

	return conf.Dispatcher.Cache
}

func ResolveDispatcherCircuitBreakerConfig(conf *config.Config) *circuitbreaker.Config {
	return &conf.Dispatcher.CircuitBreaker
}

func ResolveDispatcherMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Dispatcher.Metrics
}

func ResolveDispatcherSenderConfig(conf *config.Config, logger logging.Logger) *sender.Config {
	return &conf.Dispatcher.Sender
}
