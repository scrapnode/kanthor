//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/sdk"
	uc "github.com/scrapnode/kanthor/usecases/sdk"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

func InitializeSdk(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		sdk.New,
		infrastructure.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		database.New,
		InitializeSdkUsecase,
	)
	return nil, nil
}

func InitializeSdkUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure, db database.Database) (uc.Sdk, error) {
	wire.Build(
		uc.New,
		repos.New,
	)
	return nil, nil
}
