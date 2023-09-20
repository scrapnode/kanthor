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
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/signature"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dispatcher"
	"github.com/scrapnode/kanthor/services/portalapi"
	"github.com/scrapnode/kanthor/services/scheduler"
	"github.com/scrapnode/kanthor/services/sdkapi"
	"github.com/scrapnode/kanthor/services/storage"
	dispatcher2 "github.com/scrapnode/kanthor/usecases/dispatcher"
	"github.com/scrapnode/kanthor/usecases/portal"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
	scheduler2 "github.com/scrapnode/kanthor/usecases/scheduler"
	repos2 "github.com/scrapnode/kanthor/usecases/scheduler/repos"
	"github.com/scrapnode/kanthor/usecases/sdk"
	repos3 "github.com/scrapnode/kanthor/usecases/sdk/repos"
	storage2 "github.com/scrapnode/kanthor/usecases/storage"
	repos4 "github.com/scrapnode/kanthor/usecases/storage/repos"
)

// Injectors from wire_dispatcher.go:

func InitializeDispatcher(conf *config.Config, logger logging.Logger) (services.Service, error) {
	subscriberConfig := ResolveDispatcherSubscriberConfig(conf)
	subscriber, err := streaming.NewSubscriber(subscriberConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolveDispatcherMetricsConfig(conf)
	metrics, err := metric.New(metricConfig, logger)
	if err != nil {
		return nil, err
	}
	dispatcherDispatcher, err := InitializeDispatcherUsecase(conf, logger, metrics)
	if err != nil {
		return nil, err
	}
	service := dispatcher.New(conf, logger, subscriber, metrics, dispatcherDispatcher)
	return service, nil
}

func InitializeDispatcherUsecase(conf *config.Config, logger logging.Logger, metrics metric.Metrics) (dispatcher2.Dispatcher, error) {
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
	dispatcherDispatcher := dispatcher2.New(conf, logger, timerTimer, publisher, send, cacheCache, circuitBreaker, metrics)
	return dispatcherDispatcher, nil
}

// Injectors from wire_portal_api.go:

func InitializePortalApi(conf *config.Config, logger logging.Logger) (services.Service, error) {
	idempotencyConfig := &conf.Idempotency
	idempotencyIdempotency, err := idempotency.New(idempotencyConfig, logger)
	if err != nil {
		return nil, err
	}
	coordinatorConfig := &conf.Coordinator
	coordinatorCoordinator, err := coordinator.New(coordinatorConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolvePortalApiMetricsConfig(conf)
	metrics, err := metric.New(metricConfig, logger)
	if err != nil {
		return nil, err
	}
	authenticatorConfig := ResolvePortalApiAuthenticatorConfig(conf)
	authenticatorAuthenticator, err := authenticator.New(authenticatorConfig, logger)
	if err != nil {
		return nil, err
	}
	authorizatorConfig := ResolvePortalApiAuthorizatorConfig(conf)
	authorizatorAuthorizator, err := authorizator.New(authorizatorConfig, logger)
	if err != nil {
		return nil, err
	}
	portal, err := InitializePortalUsecase(conf, logger, metrics)
	if err != nil {
		return nil, err
	}
	service := portalapi.New(conf, logger, idempotencyIdempotency, coordinatorCoordinator, metrics, authenticatorAuthenticator, authorizatorAuthorizator, portal)
	return service, nil
}

func InitializePortalUsecase(conf *config.Config, logger logging.Logger, metrics metric.Metrics) (portal.Portal, error) {
	cryptographyConfig := &conf.Cryptography
	cryptographyCryptography, err := cryptography.New(cryptographyConfig)
	if err != nil {
		return nil, err
	}
	timerTimer := timer.New()
	cacheConfig := ResolvePortalApiCacheConfig(conf)
	cacheCache, err := cache.New(cacheConfig, logger)
	if err != nil {
		return nil, err
	}
	databaseConfig := &conf.Database
	repositories := repos.New(databaseConfig, logger, timerTimer)
	portalPortal := portal.New(conf, logger, cryptographyCryptography, metrics, timerTimer, cacheCache, repositories)
	return portalPortal, nil
}

// Injectors from wire_scheduler.go:

func InitializeScheduler(conf *config.Config, logger logging.Logger) (services.Service, error) {
	subscriberConfig := ResolveSchedulerSubscriberConfig(conf)
	subscriber, err := streaming.NewSubscriber(subscriberConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolveSchedulerMetricsConfig(conf)
	metrics, err := metric.New(metricConfig, logger)
	if err != nil {
		return nil, err
	}
	schedulerScheduler, err := InitializeSchedulerUsecase(conf, logger, metrics)
	if err != nil {
		return nil, err
	}
	service := scheduler.New(conf, logger, subscriber, metrics, schedulerScheduler)
	return service, nil
}

func InitializeSchedulerUsecase(conf *config.Config, logger logging.Logger, metrics metric.Metrics) (scheduler2.Scheduler, error) {
	timerTimer := timer.New()
	signatureSignature := signature.New()
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
	repositories := repos2.New(databaseConfig, logger, timerTimer)
	schedulerScheduler := scheduler2.New(conf, logger, timerTimer, signatureSignature, publisher, cacheCache, metrics, repositories)
	return schedulerScheduler, nil
}

// Injectors from wire_sdk_api.go:

func InitializeSdkApi(conf *config.Config, logger logging.Logger) (services.Service, error) {
	idempotencyConfig := &conf.Idempotency
	idempotencyIdempotency, err := idempotency.New(idempotencyConfig, logger)
	if err != nil {
		return nil, err
	}
	coordinatorConfig := &conf.Coordinator
	coordinatorCoordinator, err := coordinator.New(coordinatorConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolveSdkApiMetricsConfig(conf)
	metrics, err := metric.New(metricConfig, logger)
	if err != nil {
		return nil, err
	}
	authorizatorConfig := ResolveSdkApiAuthorizatorConfig(conf)
	authorizatorAuthorizator, err := authorizator.New(authorizatorConfig, logger)
	if err != nil {
		return nil, err
	}
	sdk, err := InitializeSdkUsecase(conf, logger, metrics)
	if err != nil {
		return nil, err
	}
	service := sdkapi.New(conf, logger, idempotencyIdempotency, coordinatorCoordinator, metrics, authorizatorAuthorizator, sdk)
	return service, nil
}

func InitializeSdkUsecase(conf *config.Config, logger logging.Logger, metrics metric.Metrics) (sdk.Sdk, error) {
	cryptographyConfig := &conf.Cryptography
	cryptographyCryptography, err := cryptography.New(cryptographyConfig)
	if err != nil {
		return nil, err
	}
	timerTimer := timer.New()
	cacheConfig := ResolveSdkApiCacheConfig(conf)
	cacheCache, err := cache.New(cacheConfig, logger)
	if err != nil {
		return nil, err
	}
	publisherConfig := ResolveSdkApiPublisherConfig(conf)
	publisher, err := streaming.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, err
	}
	databaseConfig := &conf.Database
	repositories := repos3.New(databaseConfig, logger, timerTimer)
	sdkSdk := sdk.New(conf, logger, cryptographyCryptography, metrics, timerTimer, cacheCache, publisher, repositories)
	return sdkSdk, nil
}

// Injectors from wire_storage.go:

func InitializeStorage(conf *config.Config, logger logging.Logger) (services.Service, error) {
	subscriberConfig := ResolveStorageSubscriberConfig(conf)
	subscriber, err := streaming.NewSubscriber(subscriberConfig, logger)
	if err != nil {
		return nil, err
	}
	metricConfig := ResolveStorageMetricsConfig(conf)
	metrics, err := metric.New(metricConfig, logger)
	if err != nil {
		return nil, err
	}
	storageStorage, err := InitializeStorageUsecase(conf, logger, metrics)
	if err != nil {
		return nil, err
	}
	service := storage.New(conf, logger, subscriber, metrics, storageStorage)
	return service, nil
}

func InitializeStorageUsecase(conf *config.Config, logger logging.Logger, metrics metric.Metrics) (storage2.Storage, error) {
	datastoreConfig := &conf.Datastore
	timerTimer := timer.New()
	repositories := repos4.New(datastoreConfig, logger, timerTimer)
	storageStorage := storage2.New(conf, logger, metrics, repositories)
	return storageStorage, nil
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

func ResolveDispatcherMetricsConfig(conf *config.Config) *metric.Config {
	return &conf.Dispatcher.Metrics
}

// wire_portal_api.go:

func ResolvePortalApiAuthenticatorConfig(conf *config.Config) *authenticator.Config {
	return &conf.PortalApi.Authenticator
}

func ResolvePortalApiAuthorizatorConfig(conf *config.Config) *authorizator.Config {
	return &conf.PortalApi.Authorizator
}

func ResolvePortalApiCacheConfig(conf *config.Config) *cache.Config {
	return &conf.PortalApi.Cache
}

func ResolvePortalApiMetricsConfig(conf *config.Config) *metric.Config {
	return &conf.PortalApi.Metrics
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

func ResolveSchedulerMetricsConfig(conf *config.Config) *metric.Config {
	return &conf.Scheduler.Metrics
}

// wire_sdk_api.go:

func ResolveSdkApiCacheConfig(conf *config.Config) *cache.Config {
	return &conf.SdkApi.Cache
}

func ResolveSdkApiPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &conf.SdkApi.Publisher
}

func ResolveSdkApiAuthorizatorConfig(conf *config.Config) *authorizator.Config {
	return &conf.SdkApi.Authorizator
}

func ResolveSdkApiMetricsConfig(conf *config.Config) *metric.Config {
	return &conf.SdkApi.Metrics
}

// wire_storage.go:

func ResolveStorageSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Storage.Subscriber
}

func ResolveStorageMetricsConfig(conf *config.Config) *metric.Config {
	return &conf.Storage.Metrics
}
