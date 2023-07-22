// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/controlplane"
	"github.com/scrapnode/kanthor/services/dataplane"
	"github.com/scrapnode/kanthor/services/dispatcher"
	"github.com/scrapnode/kanthor/services/scheduler"
	"github.com/scrapnode/kanthor/usecases"
	controlplane2 "github.com/scrapnode/kanthor/usecases/controlplane"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
	repos2 "github.com/scrapnode/kanthor/usecases/dataplane/repos"
	repos3 "github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

// Injectors from wire_controlplane.go:

func InitializeControlplane(conf *config.Config, logger logging.Logger) (services.Service, error) {
	authenticatorConfig := ResolveControlplaneAuthenticatorConfig(conf)
	authenticatorAuthenticator, err := authenticator.New(authenticatorConfig, logger)
	if err != nil {
		return nil, err
	}
	authorizatorConfig := ResolveControlplaneAuthorizatorConfig(conf)
	authorizatorAuthorizator, err := authorizator.New(authorizatorConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolveControlplaneMetricConfig(conf)
	meter := metric.New(metricConfig)
	timerTimer := timer.New()
	cacheConfig := ResolveControlplaneCacheConfig(conf)
	cacheCache, err := cache.New(cacheConfig, logger)
	if err != nil {
		return nil, err
	}
	databaseConfig := &conf.Database
	repositories := repos.New(databaseConfig, logger, timerTimer)
	controlplaneControlplane := usecases.NewControlplane(conf, logger, timerTimer, cacheCache, meter, authorizatorAuthorizator, repositories)
	service := controlplane.New(conf, logger, authenticatorAuthenticator, authorizatorAuthorizator, meter, controlplaneControlplane)
	return service, nil
}

func InitializeControlplaneUsecase(conf *config.Config, logger logging.Logger) (controlplane2.Controlplane, error) {
	timerTimer := timer.New()
	cacheConfig := ResolveControlplaneCacheConfig(conf)
	cacheCache, err := cache.New(cacheConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolveControlplaneMetricConfig(conf)
	meter := metric.New(metricConfig)
	authorizatorConfig := ResolveControlplaneAuthorizatorConfig(conf)
	authorizatorAuthorizator, err := authorizator.New(authorizatorConfig, logger)
	if err != nil {
		return nil, err
	}
	databaseConfig := &conf.Database
	repositories := repos.New(databaseConfig, logger, timerTimer)
	controlplaneControlplane := usecases.NewControlplane(conf, logger, timerTimer, cacheCache, meter, authorizatorAuthorizator, repositories)
	return controlplaneControlplane, nil
}

// Injectors from wire_dataplane.go:

func InitializeDataplane(conf *config.Config, logger logging.Logger) (services.Service, error) {
	authenticatorConfig := ResolveDataplaneAuthenticatorConfig(conf)
	authenticatorAuthenticator, err := authenticator.New(authenticatorConfig, logger)
	if err != nil {
		return nil, err
	}
	authorizatorConfig := ResolveDataplaneAuthorizatorConfig(conf)
	authorizatorAuthorizator, err := authorizator.New(authorizatorConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolveDataplaneMetricConfig(conf)
	meter := metric.New(metricConfig)
	timerTimer := timer.New()
	publisherConfig := ResolveDataplanePublisherConfig(conf)
	publisher, err := streaming.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, err
	}
	cacheConfig := ResolveDataplaneCacheConfig(conf)
	cacheCache, err := cache.New(cacheConfig, logger)
	if err != nil {
		return nil, err
	}
	databaseConfig := &conf.Database
	repositories := repos2.New(databaseConfig, logger, timerTimer)
	dataplaneDataplane := usecases.NewDataplane(conf, logger, timerTimer, publisher, cacheCache, meter, repositories)
	service := dataplane.New(conf, logger, authenticatorAuthenticator, authorizatorAuthorizator, meter, dataplaneDataplane)
	return service, nil
}

// Injectors from wire_dispatcher.go:

func InitializeDispatcher(conf *config.Config, logger logging.Logger) (services.Service, error) {
	subscriberConfig := ResolveDispatcherSubscriberConfig(conf)
	subscriber, err := streaming.NewSubscriber(subscriberConfig, logger)
	if err != nil {
		return nil, err
	}
	timerTimer := timer.New()
	publisherConfig := ResolveDispatcherPublisherConfig(conf)
	publisher, err := streaming.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, err
	}
	senderConfig := ResolveDispatcherSenderConfig(conf, logger)
	send := sender.New(senderConfig, logger)
	cacheConfig := ResolveDispatcherCacheConfig(conf)
	cacheCache, err := cache.New(cacheConfig, logger)
	if err != nil {
		return nil, err
	}
	circuitbreakerConfig := ResolveDispatcherCircuitBreakerConfig(conf)
	circuitBreaker, err := circuitbreaker.New(circuitbreakerConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolveDispatcherMetricConfig(conf)
	meter := metric.New(metricConfig)
	dispatcherDispatcher := usecases.NewDispatcher(conf, logger, timerTimer, publisher, send, cacheCache, circuitBreaker, meter)
	service := dispatcher.New(conf, logger, subscriber, dispatcherDispatcher, meter)
	return service, nil
}

// Injectors from wire_scheduler.go:

func InitializeScheduler(conf *config.Config, logger logging.Logger) (services.Service, error) {
	subscriberConfig := ResolveSchedulerSubscriberConfig(conf)
	subscriber, err := streaming.NewSubscriber(subscriberConfig, logger)
	if err != nil {
		return nil, err
	}
	timerTimer := timer.New()
	publisherConfig := ResolveSchedulerPublisherConfig(conf)
	publisher, err := streaming.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, err
	}
	cacheConfig := ResolveSchedulerCacheConfig(conf)
	cacheCache, err := cache.New(cacheConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolveSchedulerMetricConfig(conf)
	meter := metric.New(metricConfig)
	databaseConfig := &conf.Database
	repositories := repos3.New(databaseConfig, logger, timerTimer)
	schedulerScheduler := usecases.NewScheduler(conf, logger, timerTimer, publisher, cacheCache, meter, repositories)
	service := scheduler.New(conf, logger, subscriber, schedulerScheduler, meter)
	return service, nil
}

// wire_controlplane.go:

func ResolveControlplaneCacheConfig(conf *config.Config) *cache.Config {
	return &conf.Controlplane.Cache
}

func ResolveControlplaneAuthenticatorConfig(conf *config.Config) *authenticator.Config {
	return &conf.Controlplane.Authenticator
}

func ResolveControlplaneAuthorizatorConfig(conf *config.Config) *authorizator.Config {
	return &conf.Controlplane.Authorizator
}

func ResolveControlplaneMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Controlplane.Metrics
}

// wire_dataplane.go:

func ResolveDataplanePublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &conf.Dataplane.Publisher
}

func ResolveDataplaneCacheConfig(conf *config.Config) *cache.Config {
	return &conf.Dataplane.Cache
}

func ResolveDataplaneAuthenticatorConfig(conf *config.Config) *authenticator.Config {
	return &conf.Dataplane.Authenticator
}

func ResolveDataplaneAuthorizatorConfig(conf *config.Config) *authorizator.Config {
	return &conf.Dataplane.Authorizator
}

func ResolveDataplaneMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Dataplane.Metrics
}

// wire_dispatcher.go:

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

func ResolveDispatcherMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Dispatcher.Metrics
}

func ResolveDispatcherSenderConfig(conf *config.Config, logger logging.Logger) *sender.Config {
	return &conf.Dispatcher.Sender
}

// wire_scheduler.go:

func ResolveSchedulerPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &conf.Scheduler.Publisher
}

func ResolveSchedulerSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Scheduler.Subscriber
}

func ResolveSchedulerCacheConfig(conf *config.Config) *cache.Config {
	return &conf.Scheduler.Cache
}

func ResolveSchedulerMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Scheduler.Metrics
}
