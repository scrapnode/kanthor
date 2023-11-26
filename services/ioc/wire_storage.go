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
	"github.com/scrapnode/kanthor/services/storage/config"
	"github.com/scrapnode/kanthor/services/storage/entrypoint"
	"github.com/scrapnode/kanthor/services/storage/repositories"
	"github.com/scrapnode/kanthor/services/storage/usecase"
)

func Storage(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		repositories.New,
		usecase.New,
		entrypoint.Consumer,
	)
	return nil, nil
}
