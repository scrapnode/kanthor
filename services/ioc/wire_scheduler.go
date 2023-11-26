//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/scheduler/config"
	"github.com/scrapnode/kanthor/services/scheduler/entrypoint"
	"github.com/scrapnode/kanthor/services/scheduler/repositories"
	"github.com/scrapnode/kanthor/services/scheduler/usecase"
)

func Scheduler(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Consumer,
	)
	return nil, nil
}
