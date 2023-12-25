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
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/entrypoint"
	"github.com/scrapnode/kanthor/services/portal/repositories"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

func Portal(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		timer.New,
		infrastructure.New,
		database.New,
		datastore.New,
		repositories.New,
		usecase.New,
		entrypoint.Rest,
	)
	return nil, nil
}
