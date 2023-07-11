//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane"
	"github.com/scrapnode/kanthor/usecases"
)

func InitializeDataplane(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		dataplane.New,
		usecases.NewDataplane,
		timer.New,
		ResolveDataplanePublisherConfig,
		streaming.NewPublisher,
		wire.FieldsOf(new(*config.Config), "Database"),
		repositories.New,
		ResolveDataplaneCacheConfig,
		cache.New,
		ResolveDataplaneAuthenticatorConfig,
		authenticator.New,
		ResolveDataplaneMetricConfig,
		metric.New,
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

func ResolveDataplaneAuthenticatorConfig(conf *config.Config) *authenticator.Config {
	return &conf.Dataplane.Authenticator
}

func ResolveDataplaneMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Dataplane.Metrics
}
