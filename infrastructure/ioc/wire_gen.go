// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

// Injectors from wire.go:

func InitializeLogger(conf *logging.Config) (logging.Logger, error) {
	logger, err := logging.New(conf)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func InitializeStreamingPublisher(conf *streaming.PublisherConfig, logger logging.Logger) (streaming.Publisher, error) {
	publisher := streaming.NewPublisher(conf, logger)
	return publisher, nil
}

func InitializeDatabase(conf *database.Config, logger logging.Logger) (database.Database, error) {
	databaseDatabase := database.New(conf, logger)
	return databaseDatabase, nil
}

func InitializeDatastore(conf *datastore.Config, logger logging.Logger) (datastore.Datastore, error) {
	datastoreDatastore := datastore.New(conf, logger)
	return datastoreDatastore, nil
}
