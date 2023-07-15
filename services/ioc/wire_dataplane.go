//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane"
	"github.com/scrapnode/kanthor/usecases"
	"github.com/scrapnode/kanthor/usecases/dataplane/repos"
)

func InitializeDataplane(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		dataplane.New,
		usecases.NewDataplane,
		timer.New,
		ResolveDataplanePublisherConfig,
		streaming.NewPublisher,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
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
	return &conf.Dataplane.Publisher
}

func ResolveDataplaneCacheConfig(conf *config.Config) *cache.Config {
	return &conf.Dataplane.Cache
}

func ResolveDataplaneAuthenticatorConfig(conf *config.Config) *authenticator.Config {
	return &conf.Dataplane.Authenticator
}

func ResolveDataplaneMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Dataplane.Metrics
}
