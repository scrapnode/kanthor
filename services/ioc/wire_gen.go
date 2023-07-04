// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
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

// Injectors from wire.go:

func InitializeMigration(conf *config.Config, logger logging.Logger) (services.Service, error) {
	databaseConfig := ResolveDatabaseConfig(conf)
	databaseDatabase := database.New(databaseConfig, logger)
	service := migration.New(conf, logger, databaseDatabase)
	return service, nil
}

func InitializeDataplane(conf *config.Config, logger logging.Logger) (services.Service, error) {
	timerTimer := timer.New()
	publisherConfig := ResolveDataplanePublisherConfig(conf)
	publisher := streaming.NewPublisher(publisherConfig, logger)
	databaseConfig := ResolveDatabaseConfig(conf)
	repositoriesRepositories := repositories.New(databaseConfig, logger, timerTimer)
	cacheConfig := ResolveDataplaneCacheConfig(conf)
	cacheCache := cache.New(cacheConfig, logger)
	dataplaneDataplane := usecases.NewDataplane(conf, logger, timerTimer, publisher, repositoriesRepositories, cacheCache)
	metricConfig := ResolveDataplaneMetricConfig(conf)
	meter := metric.New(metricConfig)
	service := dataplane.New(conf, logger, dataplaneDataplane, meter)
	return service, nil
}

func InitializeScheduler(conf *config.Config, logger logging.Logger) (services.Service, error) {
	subscriberConfig := ResolveSchedulerSubscriberConfig(conf)
	subscriber := streaming.NewSubscriber(subscriberConfig, logger)
	timerTimer := timer.New()
	publisherConfig := ResolveSchedulerPublisherConfig(conf)
	publisher := streaming.NewPublisher(publisherConfig, logger)
	databaseConfig := ResolveDatabaseConfig(conf)
	repositoriesRepositories := repositories.New(databaseConfig, logger, timerTimer)
	cacheConfig := ResolveSchedulerCacheConfig(conf)
	cacheCache := cache.New(cacheConfig, logger)
	schedulerScheduler := usecases.NewScheduler(conf, logger, timerTimer, publisher, repositoriesRepositories, cacheCache)
	metricConfig := ResolveSchedulerMetricConfig(conf)
	meter := metric.New(metricConfig)
	service := scheduler.New(conf, logger, subscriber, schedulerScheduler, meter)
	return service, nil
}

func InitializeDispatcher(conf *config.Config, logger logging.Logger) (services.Service, error) {
	subscriberConfig := ResolveDispatcherSubscriberConfig(conf)
	subscriber := streaming.NewSubscriber(subscriberConfig, logger)
	timerTimer := timer.New()
	publisherConfig := ResolveDispatcherPublisherConfig(conf)
	publisher := streaming.NewPublisher(publisherConfig, logger)
	databaseConfig := ResolveDatabaseConfig(conf)
	repositoriesRepositories := repositories.New(databaseConfig, logger, timerTimer)
	send := ResolveDispatcherSender(conf, logger)
	cacheConfig := ResolveDispatcherCacheConfig(conf)
	cacheCache := cache.New(cacheConfig, logger)
	circuitbreakerConfig := ResolveDispatcherCircuitBreakerConfig(conf)
	circuitBreaker := circuitbreaker.New(circuitbreakerConfig, logger)
	dispatcherDispatcher := usecases.NewDispatcher(conf, logger, timerTimer, publisher, repositoriesRepositories, send, cacheCache, circuitBreaker)
	metricConfig := ResolveDispatcherMetricConfig(conf)
	meter := metric.New(metricConfig)
	service := dispatcher.New(conf, logger, subscriber, dispatcherDispatcher, meter)
	return service, nil
}

// wire.go:

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

func ResolveDataplaneMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Dataplane.Metrics
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

func ResolveSchedulerMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Scheduler.Metrics
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

func ResolveDatabaseConfig(conf *config.Config) *database.Config {
	return &conf.Database
}

func ResolveDispatcherSender(conf *config.Config, logger logging.Logger) sender.Send {
	return sender.New(&conf.Dispatcher.Sender, logger)
}
