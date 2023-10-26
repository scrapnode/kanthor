//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/scheduler"
	uc "github.com/scrapnode/kanthor/usecases/scheduler"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

func InitializeScheduler(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		scheduler.New,
		infrastructure.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		database.New,
		InitializeSchedulerUsecase,
	)
	return nil, nil
}

func InitializeSchedulerUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure, db database.Database) (uc.Scheduler, error) {
	wire.Build(
		uc.New,
		repos.New,
	)
	return nil, nil
}
