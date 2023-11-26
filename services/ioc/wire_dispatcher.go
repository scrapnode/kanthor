//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/dispatcher/config"
	"github.com/scrapnode/kanthor/services/dispatcher/entrypoint"
	"github.com/scrapnode/kanthor/services/dispatcher/usecase"
)

func Dispatcher(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		usecase.New,
		entrypoint.Consumer,
	)
	return nil, nil
}
