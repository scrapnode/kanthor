//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
	"github.com/scrapnode/kanthor/services/dispatcher"
)

func InitializeDispatcher(conf *config.Config, logger logging.Logger) (patterns.Runnable, error) {
	wire.Build(
		dispatcher.New,
		infrastructure.New,
		usecase.New,
	)
	return nil, nil
}
