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
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/controlplane"
	"github.com/scrapnode/kanthor/usecases"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
)

func InitializeControlplane(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		controlplane.New,
		usecases.NewControlplane,
		timer.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
		ResolveControlplaneCacheConfig,
		cache.New,
		ResolveControlplaneAuthenticatorConfig,
		authenticator.New,
		ResolveControlplaneMetricConfig,
		metric.New,
	)
	return nil, nil
}

func ResolveControlplaneCacheConfig(conf *config.Config) *cache.Config {
	if conf.Controlplane.Cache == nil {
		return &conf.Cache
	}

	return conf.Controlplane.Cache
}

func ResolveControlplaneAuthenticatorConfig(conf *config.Config) *authenticator.Config {
	return &conf.Controlplane.Authenticator
}

func ResolveControlplaneMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Controlplane.Metrics
}
