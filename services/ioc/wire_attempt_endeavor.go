//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/entrypoint/endeavor/executor"
	"github.com/scrapnode/kanthor/services/attempt/entrypoint/endeavor/planner"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func AttemptEndeavorPlanner(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		repositories.New,
		usecase.New,
		planner.New,
	)
	return nil, nil
}

func AttemptEndeavorExecutor(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		repositories.New,
		usecase.New,
		executor.New,
	)
	return nil, nil
}
