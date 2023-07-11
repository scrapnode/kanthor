//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/scheduler"
	"github.com/scrapnode/kanthor/usecases"
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
		repositories.New,
		ResolveSchedulerCacheConfig,
		cache.New,
		ResolveSchedulerMetricConfig,
		metric.New,
	)
	return nil, nil
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
