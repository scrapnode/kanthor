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
	"github.com/scrapnode/kanthor/services/portal/config"
	entrypoint "github.com/scrapnode/kanthor/services/portal/entrypoint/rest"
	"github.com/scrapnode/kanthor/services/portal/repositories"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

func Portal(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.New,
	)
	return nil, nil
}
