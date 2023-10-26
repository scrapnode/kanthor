//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/storage"
	uc "github.com/scrapnode/kanthor/usecases/storage"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

func InitializeStorage(conf *config.Config, logger logging.Logger) (services.Service, error) {
	wire.Build(
		storage.New,
		infrastructure.New,
		wire.FieldsOf(new(*config.Config), "Datastore"),
		datastore.New,
		InitializeStorageUsecase,
	)
	return nil, nil
}

func InitializeStorageUsecase(conf *config.Config, logger logging.Logger, infra *infrastructure.Infrastructure, ds datastore.Datastore) (uc.Storage, error) {
	wire.Build(
		uc.New,
		repos.New,
	)
	return nil, nil
}
