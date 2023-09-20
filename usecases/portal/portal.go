package portal

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type Portal interface {
	patterns.Connectable
	Account() Account
	Workspace() Workspace
	WorkspaceCredentials() WorkspaceCredentials
}

func New(
	conf *config.Config,
	logger logging.Logger,
	cryptography cryptography.Cryptography,
	metrics metric.Metrics,
	timer timer.Timer,
	cache cache.Cache,
	repos repos.Repositories,
) Portal {
	return &portal{
		conf:         conf,
		logger:       logger,
		cryptography: cryptography,
		metrics:      metrics,
		timer:        timer,
		cache:        cache,
		repos:        repos,
	}
}

type portal struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	metrics      metric.Metrics
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories

	mu                   sync.RWMutex
	account              *account
	workspace            *workspace
	workspaceCredentials *workspaceCredentials
}

func (uc *portal) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

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
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.logger.Info("disconnected")

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *portal) Account() Account {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.account == nil {
		uc.account = &account{
			conf:         uc.conf,
			logger:       uc.logger,
			cryptography: uc.cryptography,
			metrics:      uc.metrics,
			timer:        uc.timer,
			cache:        uc.cache,
			repos:        uc.repos,
		}
	}
	return uc.account
}

func (uc *portal) Workspace() Workspace {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.workspace == nil {
		uc.workspace = &workspace{
			conf:         uc.conf,
			logger:       uc.logger,
			cryptography: uc.cryptography,
			metrics:      uc.metrics,
			timer:        uc.timer,
			cache:        uc.cache,
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
			metrics:      uc.metrics,
			timer:        uc.timer,
			cache:        uc.cache,
			repos:        uc.repos,
		}
	}
	return uc.workspaceCredentials
}
