//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers"
	confprovider "github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/ioc"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/msgbroker"
)

func InitializeConfig(provider confprovider.Provider) (*config.Config, error) {
	wire.Build(config.New)
	return nil, nil
}

func GetLoggerConfig(conf *config.Config) *logging.Config {
	return conf.Logger
}

func GetMsgBrokerConfig(conf *config.Config) *msgbroker.Config {
	return conf.MsgBroker
}

func GetDatabaseConfig(conf *config.Config) *database.Config {
	return conf.Database
}

func InitializeServer(provider confprovider.Provider) (servers.Servers, error) {
	wire.Build(
		servers.New,
		InitializeConfig,
		GetLoggerConfig,
		ioc.InitializeLogger,
		GetMsgBrokerConfig,
		ioc.InitializeMsgBroker,
		GetDatabaseConfig,
		ioc.InitializeDatabase,
	)
	return nil, nil
}
