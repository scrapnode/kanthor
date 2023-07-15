//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/scheduler"
	"github.com/scrapnode/kanthor/usecases"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

func InitializeScheduler(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		scheduler.New,
		usecases.NewScheduler,
		timer.New,
		ResolveSchedulerPublisherConfig,
		streaming.NewPublisher,
		ResolveSchedulerSubscriberConfig,
		streaming.NewSubscriber,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
		ResolveSchedulerCacheConfig,
		cache.New,
		ResolveSchedulerMetricConfig,
		metric.New,
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

func ResolveSchedulerMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Scheduler.Metrics
}
