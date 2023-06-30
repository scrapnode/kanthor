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
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane"
	"github.com/scrapnode/kanthor/services/dispatcher"
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
		ResolveSchedulerSubscriberConfig,
		streaming.NewSubscriber,
		ResolveDatabaseConfig,
		repositories.New,
	)
	return nil, nil
}

func InitializeDispatcher(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		dispatcher.New,
		usecases.NewDispatcher,
		timer.New,
		ResolvePublisherConfig,
		streaming.NewPublisher,
		ResolveDispatcherSubscriberConfig,
		streaming.NewSubscriber,
		ResolveDatabaseConfig,
		repositories.New,
		ResolveSender,
	)
	return nil, nil
}

func ResolvePublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &streaming.PublisherConfig{ConnectionConfig: conf.Streaming}
}

func ResolveSchedulerSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Scheduler.Consumer
}

func ResolveDispatcherSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Dispatcher.Consumer
}

func ResolveDatabaseConfig(conf *config.Config) *database.Config {
	return &conf.Database
}

func ResolveSender(conf *config.Config, logger logging.Logger) sender.Send {
	return sender.New(&conf.Dispatcher.Sender, logger)
}
