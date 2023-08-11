//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
)

func InitializeLogger(conf *logging.Config) (logging.Logger, error) {
	wire.Build(logging.New)
	return nil, nil
}

func InitializeStreamingPublisher(conf *streaming.PublisherConfig, logger logging.Logger) (streaming.Publisher, error) {
	wire.Build(streaming.NewPublisher)
	return nil, nil
}

func InitializeDatabase(conf *database.Config, logger logging.Logger, timer timer.Timer) (database.Database, error) {
	wire.Build(database.New)
	return nil, nil
}

func InitializeDatastore(conf *datastore.Config, logger logging.Logger, timer timer.Timer) (datastore.Datastore, error) {
	wire.Build(datastore.New)
	return nil, nil
}
