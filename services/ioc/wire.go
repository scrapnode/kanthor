//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane"
	"github.com/scrapnode/kanthor/services/dispatcher"
	"github.com/scrapnode/kanthor/services/migration"
	"github.com/scrapnode/kanthor/services/scheduler"
	"github.com/scrapnode/kanthor/usecases"
)

func InitializeMigration(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		migration.New,
		ResolveDatabaseConfig,
		database.New,
	)
	return nil, nil
}

func InitializeDataplane(conf *config.Config, logger logging.Logger, meter metric.Meter) (services.Service, error) {
	wire.Build(
		dataplane.New,
		usecases.NewDataplane,
		timer.New,
		ResolveDataplanePublisherConfig,
		streaming.NewPublisher,
		ResolveDatabaseConfig,
		repositories.New,
		ResolveDataplaneCacheConfig,
		cache.New,
	)
	return nil, nil
}

func InitializeScheduler(conf *config.Config, logger logging.Logger, meter metric.Meter) (services.Service, error) {
	wire.Build(
		scheduler.New,
		usecases.NewScheduler,
		timer.New,
		ResolveSchedulerPublisherConfig,
		streaming.NewPublisher,
		ResolveSchedulerSubscriberConfig,
		streaming.NewSubscriber,
		ResolveDatabaseConfig,
		repositories.New,
		ResolveSchedulerCacheConfig,
		cache.New,
	)
	return nil, nil
}

func InitializeDispatcher(conf *config.Config, logger logging.Logger, meter metric.Meter) (services.Service, error) {
	wire.Build(
		dispatcher.New,
		usecases.NewDispatcher,
		timer.New,
		ResolveDispatcherPublisherConfig,
		streaming.NewPublisher,
		ResolveDispatcherSubscriberConfig,
		streaming.NewSubscriber,
		ResolveDatabaseConfig,
		repositories.New,
		ResolveDispatcherSender,
		ResolveDispatcherCacheConfig,
		cache.New,
		circuitbreaker.New,
		ResolveDispatcherCircuitBreakerConfig,
	)
	return nil, nil
}

func ResolveDataplanePublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	publisher := conf.Dataplane.Publisher
	if publisher.ConnectionConfig == nil {
		publisher.ConnectionConfig = &conf.Streaming
	}
	return &publisher
}

func ResolveDataplaneCacheConfig(conf *config.Config) *cache.Config {
	if conf.Dataplane.Cache == nil {
		return &conf.Cache
	}

	return conf.Dataplane.Cache
}

func ResolveSchedulerPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	publisher := conf.Scheduler.Publisher
	if publisher.ConnectionConfig == nil {
		publisher.ConnectionConfig = &conf.Streaming
	}
	return &publisher
}

func ResolveSchedulerSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	subscriber := conf.Scheduler.Subscriber
	if subscriber.ConnectionConfig == nil {
		subscriber.ConnectionConfig = &conf.Streaming
	}

	return &subscriber
}

func ResolveSchedulerCacheConfig(conf *config.Config) *cache.Config {
	if conf.Scheduler.Cache == nil {
		return &conf.Cache
	}

	return conf.Scheduler.Cache
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

func ResolveDatabaseConfig(conf *config.Config) *database.Config {
	return &conf.Database
}

func ResolveDispatcherSender(conf *config.Config, logger logging.Logger) sender.Send {
	return sender.New(&conf.Dispatcher.Sender, logger)
}
