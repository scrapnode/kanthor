//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/sdkapi"
	sdkuc "github.com/scrapnode/kanthor/usecases/sdk"
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

func InitializeSdkUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure) (sdkuc.Sdk, error) {
	wire.Build(
		sdkuc.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
		ResolveSdkApiPublisherConfig,
		streaming.NewPublisher,
	)
	return nil, nil
}

func ResolveSdkApiPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &conf.SdkApi.Publisher
}
