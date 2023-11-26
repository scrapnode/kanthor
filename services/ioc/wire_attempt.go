//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/entrypoint"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func AttemptTriggerPlanner(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		database.New,
		datastore.New,
		repositories.New,
		usecase.New,
		entrypoint.TriggerPlanner,
	)
	return nil, nil
}

func AttemptTriggerExecutor(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		database.New,
		datastore.New,
		repositories.New,
		usecase.New,
		entrypoint.TriggerExecutor,
	)
	return nil, nil
}

func AttemptEndeavorPlanner(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		database.New,
		datastore.New,
		repositories.New,
		usecase.New,
		entrypoint.EndeavorPlanner,
	)
	return nil, nil
}

func AttemptEndeavorExecutor(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		database.New,
		datastore.New,
		repositories.New,
		usecase.New,
		entrypoint.EndeavorExecutor,
	)
	return nil, nil
}
