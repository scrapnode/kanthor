//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/msgbroker"
)

func InitializeLogger(conf config.Provider) (logging.Logger, error) {
	wire.Build(logging.New)
	return nil, nil
}

func InitializeMsgBroker(conf config.Provider) (msgbroker.MsgBroker, error) {
	wire.Build(msgbroker.New, InitializeLogger)
	return nil, nil
}

func InitializeDatabase(conf config.Provider) (database.Database, error) {
	wire.Build(database.New, InitializeLogger)
	return nil, nil
}
