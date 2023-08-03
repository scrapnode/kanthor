// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dispatcher"
	"github.com/scrapnode/kanthor/services/scheduler"
	"github.com/scrapnode/kanthor/services/sdkapi"
	dispatcher2 "github.com/scrapnode/kanthor/usecases/dispatcher"
	"github.com/scrapnode/kanthor/usecases/portal"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
	scheduler2 "github.com/scrapnode/kanthor/usecases/scheduler"
	repos2 "github.com/scrapnode/kanthor/usecases/scheduler/repos"
	"github.com/scrapnode/kanthor/usecases/sdk"
	repos3 "github.com/scrapnode/kanthor/usecases/sdk/repos"
)

// Injectors from wire_dispatcher.go:

func InitializeDispatcher(conf *config.Config, logger logging.Logger) (services.Service, error) {
	subscriberConfig := ResolveDispatcherSubscriberConfig(conf)
	subscriber, err := streaming.NewSubscriber(subscriberConfig, logger)
	if err != nil {
		return nil, err
	}
	dispatcherDispatcher, err := InitializeDispatcherUsecase(conf, logger)
	if err != nil {
		return nil, err
	}
	service := dispatcher.New(conf, logger, subscriber, dispatcherDispatcher)
	return service, nil
}

func InitializeDispatcherUsecase(conf *config.Config, logger logging.Logger) (dispatcher2.Dispatcher, error) {
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
	dispatcherDispatcher := dispatcher2.New(conf, logger, timerTimer, publisher, send, cacheCache, circuitBreaker)
	return dispatcherDispatcher, nil
}

// Injectors from wire_portal_api.go:

func InitializePortalUsecase(conf *config.Config, logger logging.Logger) (portal.Portal, error) {
	cryptographyConfig := &conf.Cryptography
	cryptographyCryptography, err := cryptography.New(cryptographyConfig)
	if err != nil {
		return nil, err
	}
	timerTimer := timer.New()
	cacheConfig := ResolvePortalCacheConfig(conf)
	cacheCache, err := cache.New(cacheConfig, logger)
	if err != nil {
		return nil, err
	}
	databaseConfig := &conf.Database
	repositories := repos.New(databaseConfig, logger)
	portalPortal := portal.New(conf, logger, cryptographyCryptography, timerTimer, cacheCache, repositories)
	return portalPortal, nil
}

// Injectors from wire_scheduler.go:

func InitializeScheduler(conf *config.Config, logger logging.Logger) (services.Service, error) {
	subscriberConfig := ResolveSchedulerSubscriberConfig(conf)
	subscriber, err := streaming.NewSubscriber(subscriberConfig, logger)
	if err != nil {
		return nil, err
	}
	schedulerScheduler, err := InitializeSchedulerUsecase(conf, logger)
	if err != nil {
		return nil, err
	}
	service := scheduler.New(conf, logger, subscriber, schedulerScheduler)
	return service, nil
}

func InitializeSchedulerUsecase(conf *config.Config, logger logging.Logger) (scheduler2.Scheduler, error) {
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
	databaseConfig := &conf.Database
	repositories := repos2.New(databaseConfig, logger)
	schedulerScheduler := scheduler2.New(conf, logger, timerTimer, publisher, cacheCache, repositories)
	return schedulerScheduler, nil
}

// Injectors from wire_sdk_api.go:

func InitializeSdkApi(conf *config.Config, logger logging.Logger) (services.Service, error) {
	validatorValidator := validator.New()
	authorizatorConfig := ResolveSdkAuthorizatorConfig(conf)
	authorizatorAuthorizator, err := authorizator.New(authorizatorConfig, logger)
	if err != nil {
		return nil, err
	}
	sdk, err := InitializeSdkUsecase(conf, logger)
	if err != nil {
		return nil, err
	}
	service := sdkapi.New(conf, logger, validatorValidator, authorizatorAuthorizator, sdk)
	return service, nil
}

func InitializeSdkUsecase(conf *config.Config, logger logging.Logger) (sdk.Sdk, error) {
	cryptographyConfig := &conf.Cryptography
	cryptographyCryptography, err := cryptography.New(cryptographyConfig)
	if err != nil {
		return nil, err
	}
	timerTimer := timer.New()
	cacheConfig := ResolveSdkCacheConfig(conf)
	cacheCache, err := cache.New(cacheConfig, logger)
	if err != nil {
		return nil, err
	}
	databaseConfig := &conf.Database
	repositories := repos3.New(databaseConfig, logger)
	sdkSdk := sdk.New(conf, logger, cryptographyCryptography, timerTimer, cacheCache, repositories)
	return sdkSdk, nil
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

func ResolveDispatcherSenderConfig(conf *config.Config, logger logging.Logger) *sender.Config {
	return &conf.Dispatcher.Sender
}

// wire_portal_api.go:

func ResolvePortalCacheConfig(conf *config.Config) *cache.Config {
	return &conf.PortalApi.Cache
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

// wire_sdk_api.go:

func ResolveSdkCacheConfig(conf *config.Config) *cache.Config {
	return &conf.SdkApi.Cache
}

func ResolveSdkAuthorizatorConfig(conf *config.Config) *authorizator.Config {
	return &conf.SdkApi.Authorizator
}
