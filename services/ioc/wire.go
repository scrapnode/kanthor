//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane"
	"github.com/scrapnode/kanthor/services/migration"
	"github.com/scrapnode/kanthor/services/scheduler"
	"github.com/scrapnode/kanthor/usecases"
)

func InitializeMigration(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		migration.New,
		ResolveDatabaseConfig,
		database.New,
	)
	return nil, nil
}

func InitializeDataplane(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		dataplane.New,
		usecases.NewDataplane,
		timer.New,
		ResolvePublisherConfig,
		streaming.NewPublisher,
		ResolveDatabaseConfig,
		repositories.New,
	)
	return nil, nil
}

func InitializeScheduler(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		scheduler.New,
		usecases.NewScheduler,
		timer.New,
		ResolvePublisherConfig,
		streaming.NewPublisher,
		ResolveSubscriberConfig,
		streaming.NewSubscriber,
		ResolveDatabaseConfig,
		repositories.New,
	)
	return nil, nil
}

func ResolvePublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &streaming.PublisherConfig{ConnectionConfig: conf.Streaming}
}

func ResolveSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Scheduler.Consumer
}

func ResolveDatabaseConfig(conf *config.Config) *database.Config {
	return &conf.Database
}
