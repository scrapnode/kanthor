//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/portalapi"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

func InitializePortalApi(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		portalapi.New,
		infrastructure.New,
		ResolvePortalApiAuthenticatorConfig,
		authenticator.New,
		InitializePortalUsecase,
	)
	return nil, nil
}
func InitializePortalUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure) (portaluc.Portal, error) {
	wire.Build(
		portaluc.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
	)
	return nil, nil
}

func ResolvePortalApiAuthenticatorConfig(conf *config.Config) *authenticator.Config {
	return &conf.PortalApi.Authenticator
}
