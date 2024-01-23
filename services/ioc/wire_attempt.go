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

func AttemptCronjob(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Cronjob,
	)
	return nil, nil
}

func AttemptConsumer(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Consumer,
	)
	return nil, nil
}
