//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/services/attempt/repos"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
	"github.com/scrapnode/kanthor/services/scheduler"
)

func InitializeScheduler(conf *config.Config, logger logging.Logger) (patterns.Runnable, error) {
	wire.Build(
		scheduler.New,
		infrastructure.New,
		wire.FieldsOf(new(*config.Config), "Database"),
		database.New,
		repos.New,
		usecase.New,
	)
	return nil, nil
}
