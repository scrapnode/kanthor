package portal

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
	"sync"
)

type Portal interface {
	patterns.Connectable
	Workspace() Workspace
	WorkspaceCredentials() WorkspaceCredentials
}

func New(
	conf *config.Config,
	logger logging.Logger,
	cryptography cryptography.Cryptography,
	timer timer.Timer,
	cache cache.Cache,
	meter metric.Meter,
	repos repos.Repositories,
) Portal {
	return &portal{
		conf:         conf,
		logger:       logger,
		cryptography: cryptography,
		timer:        timer,
		cache:        cache,
		meter:        meter,
		repos:        repos,
	}
}

type portal struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	meter        metric.Meter
	repos        repos.Repositories

	mu                   sync.RWMutex
	workspace            *workspace
	workspaceCredentials *workspaceCredentials
}

func (usecase *portal) Connect(ctx context.Context) error {
	if err := usecase.cache.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.repos.Connect(ctx); err != nil {
		return err
	}

	usecase.logger.Info("connected")
	return nil
}

func (usecase *portal) Disconnect(ctx context.Context) error {
	usecase.logger.Info("disconnected")

	if err := usecase.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (usecase *portal) Workspace() Workspace {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.workspace == nil {
		usecase.workspace = &workspace{
			conf:         usecase.conf,
			logger:       usecase.logger,
			cryptography: usecase.cryptography,
			timer:        usecase.timer,
			cache:        usecase.cache,
			meter:        usecase.meter,
			repos:        usecase.repos,
		}
	}

	return usecase.workspace
}

func (usecase *portal) WorkspaceCredentials() WorkspaceCredentials {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.workspaceCredentials == nil {
		usecase.workspaceCredentials = &workspaceCredentials{
			conf:         usecase.conf,
			logger:       usecase.logger,
			cryptography: usecase.cryptography,
			timer:        usecase.timer,
			cache:        usecase.cache,
			meter:        usecase.meter,
			repos:        usecase.repos,
		}
	}

	return usecase.workspaceCredentials
}
