//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/infrastructure/signature"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/scheduler"
	scheduleruc "github.com/scrapnode/kanthor/usecases/scheduler"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

func InitializeScheduler(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		scheduler.New,
		ResolveSchedulerSubscriberConfig,
		streaming.NewSubscriber,
		ResolveSchedulerMetricsConfig,
		metrics.New,
		InitializeSchedulerUsecase,
	)
	return nil, nil
}

func InitializeSchedulerUsecase(conf *config.Config, logger logging.Logger, metrics metrics.Metrics) (scheduleruc.Scheduler, error) {
	wire.Build(
		scheduleruc.New,
		timer.New,
		signature.New,
		ResolveSchedulerPublisherConfig,
		streaming.NewPublisher,
		ResolveSchedulerCacheConfig,
		cache.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
	)
	return nil, nil
}

func ResolveSchedulerPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &conf.Scheduler.Publisher
}

func ResolveSchedulerSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Scheduler.Subscriber
}

func ResolveSchedulerCacheConfig(conf *config.Config) *cache.Config {
	return &conf.Scheduler.Cache
}

func ResolveSchedulerMetricsConfig(conf *config.Config) *metrics.Config {
	return &conf.Scheduler.Metrics
}
