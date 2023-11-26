// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/entrypoint"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
	config2 "github.com/scrapnode/kanthor/services/dispatcher/config"
	entrypoint2 "github.com/scrapnode/kanthor/services/dispatcher/entrypoint"
	usecase2 "github.com/scrapnode/kanthor/services/dispatcher/usecase"
	config3 "github.com/scrapnode/kanthor/services/portal/config"
	entrypoint3 "github.com/scrapnode/kanthor/services/portal/entrypoint"
	repositories2 "github.com/scrapnode/kanthor/services/portal/repositories"
	usecase3 "github.com/scrapnode/kanthor/services/portal/usecase"
	config4 "github.com/scrapnode/kanthor/services/scheduler/config"
	entrypoint4 "github.com/scrapnode/kanthor/services/scheduler/entrypoint"
	repositories3 "github.com/scrapnode/kanthor/services/scheduler/repositories"
	usecase4 "github.com/scrapnode/kanthor/services/scheduler/usecase"
	config5 "github.com/scrapnode/kanthor/services/sdk/config"
	entrypoint5 "github.com/scrapnode/kanthor/services/sdk/entrypoint"
	repositories4 "github.com/scrapnode/kanthor/services/sdk/repositories"
	usecase5 "github.com/scrapnode/kanthor/services/sdk/usecase"
	config6 "github.com/scrapnode/kanthor/services/storage/config"
	entrypoint6 "github.com/scrapnode/kanthor/services/storage/entrypoint"
	repositories5 "github.com/scrapnode/kanthor/services/storage/repositories"
	usecase6 "github.com/scrapnode/kanthor/services/storage/usecase"
)

// Injectors from wire_attempt.go:

func AttemptTriggerPlanner(provider configuration.Provider) (patterns.Runnable, error) {
	configConfig, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	databaseDatabase, err := database.New(provider)
	if err != nil {
		return nil, err
	}
	datastoreDatastore, err := datastore.New(provider)
	if err != nil {
		return nil, err
	}
	repositoriesRepositories := repositories.New(logger, databaseDatabase, datastoreDatastore)
	attempt := usecase.New(configConfig, logger, infrastructureInfrastructure, repositoriesRepositories)
	runnable := entrypoint.TriggerPlanner(configConfig, logger, infrastructureInfrastructure, databaseDatabase, datastoreDatastore, attempt)
	return runnable, nil
}

func AttemptTriggerExecutor(provider configuration.Provider) (patterns.Runnable, error) {
	configConfig, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	databaseDatabase, err := database.New(provider)
	if err != nil {
		return nil, err
	}
	datastoreDatastore, err := datastore.New(provider)
	if err != nil {
		return nil, err
	}
	repositoriesRepositories := repositories.New(logger, databaseDatabase, datastoreDatastore)
	attempt := usecase.New(configConfig, logger, infrastructureInfrastructure, repositoriesRepositories)
	runnable := entrypoint.TriggerExecutor(configConfig, logger, infrastructureInfrastructure, databaseDatabase, datastoreDatastore, attempt)
	return runnable, nil
}

func AttemptTriggerCli(provider configuration.Provider) (patterns.CommandLine, error) {
	configConfig, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	databaseDatabase, err := database.New(provider)
	if err != nil {
		return nil, err
	}
	datastoreDatastore, err := datastore.New(provider)
	if err != nil {
		return nil, err
	}
	repositoriesRepositories := repositories.New(logger, databaseDatabase, datastoreDatastore)
	attempt := usecase.New(configConfig, logger, infrastructureInfrastructure, repositoriesRepositories)
	commandLine := entrypoint.TriggerCli(configConfig, logger, infrastructureInfrastructure, databaseDatabase, datastoreDatastore, attempt)
	return commandLine, nil
}

func AttemptEndeavorPlanner(provider configuration.Provider) (patterns.Runnable, error) {
	configConfig, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	databaseDatabase, err := database.New(provider)
	if err != nil {
		return nil, err
	}
	datastoreDatastore, err := datastore.New(provider)
	if err != nil {
		return nil, err
	}
	repositoriesRepositories := repositories.New(logger, databaseDatabase, datastoreDatastore)
	attempt := usecase.New(configConfig, logger, infrastructureInfrastructure, repositoriesRepositories)
	runnable := entrypoint.EndeavorPlanner(configConfig, logger, infrastructureInfrastructure, databaseDatabase, datastoreDatastore, attempt)
	return runnable, nil
}

func AttemptEndeavorExecutor(provider configuration.Provider) (patterns.Runnable, error) {
	configConfig, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	databaseDatabase, err := database.New(provider)
	if err != nil {
		return nil, err
	}
	datastoreDatastore, err := datastore.New(provider)
	if err != nil {
		return nil, err
	}
	repositoriesRepositories := repositories.New(logger, databaseDatabase, datastoreDatastore)
	attempt := usecase.New(configConfig, logger, infrastructureInfrastructure, repositoriesRepositories)
	runnable := entrypoint.EndeavorExecutor(configConfig, logger, infrastructureInfrastructure, databaseDatabase, datastoreDatastore, attempt)
	return runnable, nil
}

// Injectors from wire_dispatcher.go:

func Dispatcher(provider configuration.Provider) (patterns.Runnable, error) {
	configConfig, err := config2.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	dispatcher := usecase2.New(configConfig, logger, infrastructureInfrastructure)
	runnable := entrypoint2.Consumer(configConfig, logger, infrastructureInfrastructure, dispatcher)
	return runnable, nil
}

// Injectors from wire_portal.go:

func Portal(provider configuration.Provider) (patterns.Runnable, error) {
	configConfig, err := config3.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	databaseDatabase, err := database.New(provider)
	if err != nil {
		return nil, err
	}
	repositoriesRepositories := repositories2.New(logger, databaseDatabase)
	portal := usecase3.New(configConfig, logger, infrastructureInfrastructure, repositoriesRepositories)
	runnable := entrypoint3.Rest(configConfig, logger, infrastructureInfrastructure, databaseDatabase, portal)
	return runnable, nil
}

// Injectors from wire_scheduler.go:

func Scheduler(provider configuration.Provider) (patterns.Runnable, error) {
	configConfig, err := config4.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	databaseDatabase, err := database.New(provider)
	if err != nil {
		return nil, err
	}
	repositoriesRepositories := repositories3.New(logger, databaseDatabase)
	scheduler := usecase4.New(configConfig, logger, infrastructureInfrastructure, repositoriesRepositories)
	runnable := entrypoint4.Consumer(configConfig, logger, infrastructureInfrastructure, databaseDatabase, scheduler)
	return runnable, nil
}

// Injectors from wire_sdk.go:

func Sdk(provider configuration.Provider) (patterns.Runnable, error) {
	configConfig, err := config5.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	databaseDatabase, err := database.New(provider)
	if err != nil {
		return nil, err
	}
	repositoriesRepositories := repositories4.New(logger, databaseDatabase)
	sdk := usecase5.New(configConfig, logger, infrastructureInfrastructure, repositoriesRepositories)
	runnable := entrypoint5.Rest(configConfig, logger, infrastructureInfrastructure, databaseDatabase, sdk)
	return runnable, nil
}

// Injectors from wire_storage.go:

func Storage(provider configuration.Provider) (patterns.Runnable, error) {
	configConfig, err := config6.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	infrastructureInfrastructure, err := infrastructure.New(provider)
	if err != nil {
		return nil, err
	}
	datastoreDatastore, err := datastore.New(provider)
	if err != nil {
		return nil, err
	}
	repositoriesRepositories := repositories5.New(logger, datastoreDatastore)
	storage := usecase6.New(configConfig, logger, infrastructureInfrastructure, repositoriesRepositories)
	runnable := entrypoint6.Consumer(configConfig, logger, infrastructureInfrastructure, datastoreDatastore, storage)
	return runnable, nil
}
