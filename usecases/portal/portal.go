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

func (uc *portal) Connect(ctx context.Context) error {
	if err := uc.cache.Connect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	uc.logger.Info("connected")
	return nil
}

func (uc *portal) Disconnect(ctx context.Context) error {
	uc.logger.Info("disconnected")

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *portal) Workspace() Workspace {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.workspace == nil {
		uc.workspace = &workspace{
			conf:         uc.conf,
			logger:       uc.logger,
			cryptography: uc.cryptography,
			timer:        uc.timer,
			cache:        uc.cache,
			meter:        uc.meter,
			repos:        uc.repos,
		}
	}
	return uc.workspace
}

func (uc *portal) WorkspaceCredentials() WorkspaceCredentials {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.workspaceCredentials == nil {
		uc.workspaceCredentials = &workspaceCredentials{
			conf:         uc.conf,
			logger:       uc.logger,
			cryptography: uc.cryptography,
			timer:        uc.timer,
			cache:        uc.cache,
			meter:        uc.meter,
			repos:        uc.repos,
		}
	}
	return uc.workspaceCredentials
}
