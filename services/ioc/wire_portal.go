//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

func InitializePortalUsecase(conf *config.Config, logger logging.Logger) (portaluc.Portal, error) {
	wire.Build(
		usecases.NewPortal,
		wire.FieldsOf(new(*config.Config), "Cryptography"),
		cryptography.New,
		timer.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
		ResolvePortalCacheConfig,
		cache.New,
		ResolvePortalMetricConfig,
		metric.New,
	)
	return nil, nil
}

func ResolvePortalCacheConfig(conf *config.Config) *cache.Config {
	return &conf.Portal.Cache
}

func ResolvePortalMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Portal.Metrics
}
