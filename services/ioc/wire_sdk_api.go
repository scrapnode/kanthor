//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/sdkapi"
	uc "github.com/scrapnode/kanthor/usecases/sdk"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

func InitializeSdkApi(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		sdkapi.New,
		infrastructure.New,
		InitializeSdkUsecase,
	)
	return nil, nil
}

func InitializeSdkUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure) (uc.Sdk, error) {
	wire.Build(
		uc.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
	)
	return nil, nil
}
