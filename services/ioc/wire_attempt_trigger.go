//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/services/attempt/repos"
	"github.com/scrapnode/kanthor/services/attempt/trigger"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func InitializeAttemptTriggerPlanner(conf *config.Config, logger logging.Logger) (patterns.Runnable, error) {
	wire.Build(
		trigger.NewPlanner,
		infrastructure.New,
		wire.FieldsOf(new(*config.Config), "Datastore"),
		datastore.New,
		InitializeAttemptUsecase,
	)
	return nil, nil
}

func InitializeAttemptTriggerExecutor(conf *config.Config, logger logging.Logger) (patterns.Runnable, error) {
	wire.Build(
		trigger.NewExecutor,
		infrastructure.New,
		wire.FieldsOf(new(*config.Config), "Datastore"),
		datastore.New,
		InitializeAttemptUsecase,
	)
	return nil, nil
}

func InitializeAttemptUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure, ds datastore.Datastore) (uc.Attempt, error) {
	wire.Build(
		usecase.New,
		repos.New,
	)
	return nil, nil
}
