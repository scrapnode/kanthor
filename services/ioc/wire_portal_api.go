//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

func InitializePortalUsecase(conf *config.Config, logger logging.Logger) (portaluc.Portal, error) {
	wire.Build(
		portaluc.New,
		wire.FieldsOf(new(*config.Config), "Cryptography"),
		cryptography.New,
		timer.New,
		ResolvePortalCacheConfig,
		cache.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
	)
	return nil, nil
}

func ResolvePortalCacheConfig(conf *config.Config) *cache.Config {
	return &conf.PortalApi.Cache
}
