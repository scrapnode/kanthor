//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/controlplane"
	"github.com/scrapnode/kanthor/usecases"
	controlplaneuc "github.com/scrapnode/kanthor/usecases/controlplane"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
)

func InitializeControlplane(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		controlplane.New,
		usecases.NewControlplane,
		wire.FieldsOf(new(*config.Config), "Symmetric"),
		cryptography.NewSymmetric,
		timer.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
		ResolveControlplaneCacheConfig,
		cache.New,
		ResolveControlplaneAuthenticatorConfig,
		authenticator.New,
		ResolveControlplaneAuthorizatorConfig,
		authorizator.New,
		ResolveControlplaneMetricConfig,
		metric.New,
	)
	return nil, nil
}

func InitializeControlplaneUsecase(conf *config.Config, logger logging.Logger) (controlplaneuc.Controlplane, error) {
	wire.Build(
		usecases.NewControlplane,
		wire.FieldsOf(new(*config.Config), "Symmetric"),
		cryptography.NewSymmetric,
		timer.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
		ResolveControlplaneCacheConfig,
		cache.New,
		ResolveControlplaneAuthorizatorConfig,
		authorizator.New,
		ResolveControlplaneMetricConfig,
		metric.New,
	)
	return nil, nil
}

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
