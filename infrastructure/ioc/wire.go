//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/msgbroker"
)

func InitializeLogger(conf *logging.Config) (logging.Logger, error) {
	wire.Build(logging.New)
	return nil, nil
}

func InitializeMsgBroker(logger logging.Logger, conf *msgbroker.Config) (msgbroker.MsgBroker, error) {
	wire.Build(msgbroker.New)
	return nil, nil
}

func InitializeDatabase(logger logging.Logger, conf *database.Config) (database.Database, error) {
	wire.Build(database.New)
	return nil, nil
}
