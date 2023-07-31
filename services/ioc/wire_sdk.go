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
	sdkuc "github.com/scrapnode/kanthor/usecases/sdk"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

func InitializeSdkUsecase(conf *config.Config, logger logging.Logger) (sdkuc.Sdk, error) {
	wire.Build(
		sdkuc.New,
		wire.FieldsOf(new(*config.Config), "Cryptography"),
		cryptography.New,
		timer.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		repos.New,
		ResolveSdkCacheConfig,
		cache.New,
		ResolveSdkMetricConfig,
		metric.New,
	)
	return nil, nil
}

func ResolveSdkCacheConfig(conf *config.Config) *cache.Config {
	return &conf.Sdk.Cache
}

func ResolveSdkMetricConfig(conf *config.Config) *metric.Config {
	return &conf.Sdk.Metrics
}
