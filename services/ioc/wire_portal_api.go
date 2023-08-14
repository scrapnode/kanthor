//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/portalapi"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

func InitializePortalApi(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		portalapi.New,
		validator.New,
		wire.FieldsOf(new(*config.Config), "Idempotency"),
		idempotency.New,
		wire.FieldsOf(new(*config.Config), "Coordinator"),
		coordinator.New,
		ResolvePortalAuthenticatorConfig,
		authenticator.New,
		ResolvePortalAuthorizatorConfig,
		authorizator.New,
		InitializePortalUsecase,
	)
	return nil, nil
}
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

func ResolvePortalAuthenticatorConfig(conf *config.Config) *authenticator.Config {
	return &conf.PortalApi.Authenticator
}

func ResolvePortalAuthorizatorConfig(conf *config.Config) *authorizator.Config {
	return &conf.PortalApi.Authorizator
}

func ResolvePortalCacheConfig(conf *config.Config) *cache.Config {
	return &conf.PortalApi.Cache
}
