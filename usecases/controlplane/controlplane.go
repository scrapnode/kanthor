package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
	"sync"
)

type Controlplane interface {
	patterns.Connectable
	Project() Project
	Workspace() Workspace
	Application() Application
}

func New(
	conf *config.Config,
	logger logging.Logger,
	symmetric cryptography.Symmetric,
	timer timer.Timer,
	cache cache.Cache,
	meter metric.Meter,
	authorizator authorizator.Authorizator,
	repos repos.Repositories,
) Controlplane {
	return &controlplane{
		conf:         conf,
		logger:       logger,
		symmetric:    symmetric,
		timer:        timer,
		cache:        cache,
		meter:        meter,
		authorizator: authorizator,
		repos:        repos,
	}
}

type controlplane struct {
	conf         *config.Config
	logger       logging.Logger
	symmetric    cryptography.Symmetric
	timer        timer.Timer
	cache        cache.Cache
	meter        metric.Meter
	authorizator authorizator.Authorizator
	repos        repos.Repositories

	mu          sync.RWMutex
	workspace   *workspace
	project     *project
	application *application
}

func (usecase *controlplane) Connect(ctx context.Context) error {
	if err := usecase.cache.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.repos.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.authorizator.Connect(ctx); err != nil {
		return err
	}

	usecase.logger.Info("connected")
	return nil
}

func (usecase *controlplane) Disconnect(ctx context.Context) error {
	usecase.logger.Info("disconnected")

	if err := usecase.authorizator.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (usecase *controlplane) Project() Project {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.project == nil {
		usecase.project = &project{
			conf:   usecase.conf,
			logger: usecase.logger,
			timer:  usecase.timer,
			cache:  usecase.cache,
			meter:  usecase.meter,
			repos:  usecase.repos,
		}
	}

	return usecase.project
}

func (usecase *controlplane) Application() Application {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.application == nil {
		usecase.application = &application{
			conf:   usecase.conf,
			logger: usecase.logger,
			timer:  usecase.timer,
			cache:  usecase.cache,
			meter:  usecase.meter,
			repos:  usecase.repos,
		}
	}

	return usecase.application
}

func (usecase *controlplane) Workspace() Workspace {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.workspace == nil {
		usecase.workspace = &workspace{
			conf:   usecase.conf,
			logger: usecase.logger,
			timer:  usecase.timer,
			cache:  usecase.cache,
			meter:  usecase.meter,
			repos:  usecase.repos,
		}
	}

	return usecase.workspace
}
