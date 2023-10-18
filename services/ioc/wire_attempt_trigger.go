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
	"github.com/scrapnode/kanthor/services/attempt/trigger"
	uc "github.com/scrapnode/kanthor/usecases/attempt"
	"github.com/scrapnode/kanthor/usecases/attempt/repos"
)

func InitializeAttemptTriggerPlanner(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		trigger.NewPlanner,
		infrastructure.New,
		InitializeAttemptUsecase,
	)
	return nil, nil
}

func InitializeAttemptTriggerExecutor(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		trigger.NewExecutor,
		ResolveAttemptTriggerExcutorSubscriberConfig,
		streaming.NewSubscriber,
		infrastructure.New,
		InitializeAttemptUsecase,
	)
	return nil, nil
}

func InitializeAttemptUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure) (uc.Attempt, error) {
	wire.Build(
		uc.New,
		ResolveAttemptTriggerPublisherConfig,
		streaming.NewPublisher,
		wire.FieldsOf(new(*config.Config), "Datastore"),
		repos.New,
	)
	return nil, nil
}

func ResolveAttemptTriggerExcutorSubscriberConfig(conf *config.Config) *streaming.SubscriberConfig {
	return &conf.Attempt.Trigger.Executor.Subscriber
}

func ResolveAttemptTriggerPublisherConfig(conf *config.Config) *streaming.PublisherConfig {
	return &conf.Attempt.Trigger.Publisher
}
