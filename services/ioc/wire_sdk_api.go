//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/sdkapi"
	sdkuc "github.com/scrapnode/kanthor/usecases/sdk"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

func InitializeSdkApi(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		sdkapi.New,
		validator.New,
		wire.FieldsOf(new(*config.Config), "Idempotency"),
		idempotency.New,
		wire.FieldsOf(new(*config.Config), "Coordinator"),
		coordinator.New,
		ResolveSdkApiMetricsConfig,
		metric.New,
		ResolveSdkApiAuthorizatorConfig,
		authorizator.New,
		InitializeSdkUsecase,
	)
	return nil, nil
}

func InitializeSdkUsecase(conf *config.Config, logger logging.Logger, metrics metric.Metrics) (sdkuc.Sdk, error) {
	wire.Build(
		sdkuc.New,
		wire.FieldsOf(new(*config.Config), "Cryptography"),
		cryptography.New,
		timer.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
		ResolveSdkApiCacheConfig,
		cache.New,
		ResolveSdkApiPublisherConfig,
		streaming.NewPublisher,
	)
	return nil, nil
}

func ResolveSdkApiCacheConfig(conf *config.Config) *cache.Config {
	return &conf.SdkApi.Cache
}

func ResolveSdkApiPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &conf.SdkApi.Publisher
}

func ResolveSdkApiAuthorizatorConfig(conf *config.Config) *authorizator.Config {
	return &conf.SdkApi.Authorizator
}

func ResolveSdkApiMetricsConfig(conf *config.Config) *metric.Config {
	return &conf.SdkApi.Metrics
}
