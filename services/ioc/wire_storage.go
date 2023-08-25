//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/storage"
	storageuc "github.com/scrapnode/kanthor/usecases/storage"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

func InitializeStorage(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		storage.New,
		ResolveStorageSubscriberConfig,
		streaming.NewSubscriber,
		ResolveStorageMetricsConfig,
		metrics.New,
		InitializeStorageUsecase,
	)
	return nil, nil
}

func InitializeStorageUsecase(conf *config.Config, logger logging.Logger, metrics metrics.Metrics) (storageuc.Storage, error) {
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

func ResolveStorageMetricsConfig(conf *config.Config) *metrics.Config {
	return &conf.Storage.Metrics
}
