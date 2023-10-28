//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/portal"
	uc "github.com/scrapnode/kanthor/usecases/portal"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

func InitializePortal(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		portal.New,
		infrastructure.New,
		ResolvePortalAuthenticatorConfig,
		authenticator.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		database.New,
		InitializePortalUsecase,
	)
	return nil, nil
}
func InitializePortalUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure, db database.Database) (uc.Portal, error) {
	wire.Build(
		uc.New,
		repos.New,
	)
	return nil, nil
}

func ResolvePortalAuthenticatorConfig(conf *config.Config) *authenticator.Config {
	return &conf.Portal.Authenticator
}
