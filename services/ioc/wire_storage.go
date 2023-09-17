//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/validation"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/storage"
	storageuc "github.com/scrapnode/kanthor/usecases/storage"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

func InitializeStorage(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		storage.New,
		validation.New,
		ResolveStorageSubscriberConfig,
		streaming.NewSubscriber,
		ResolveStorageMetricsConfig,
		metric.New,
		InitializeStorageUsecase,
	)
	return nil, nil
}

func InitializeStorageUsecase(conf *config.Config, logger logging.Logger, validator validation.Validator, metrics metric.Metrics) (storageuc.Storage, error) {
	wire.Build(
		storageuc.New,
		wire.FieldsOf(new(*config.Config), "Datastore"),
		timer.New,
		repos.New,
	)
	return nil, nil
}

func ResolveStorageSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Storage.Subscriber
}

func ResolveStorageMetricsConfig(conf *config.Config) *metric.Config {
	return &conf.Storage.Metrics
}
