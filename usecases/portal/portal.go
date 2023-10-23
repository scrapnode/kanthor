package portal

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
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
	infra *infrastructure.Infrastructure,
	repos repos.Repositories,
) Portal {
	logger = logger.With("usecase", "portal")

	return &portal{
		conf:   conf,
		logger: logger,
		infra:  infra,
		repos:  repos,
	}
}

type portal struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories

	account              *account
	workspace            *workspace
	workspaceCredentials *workspaceCredentials

	mu     sync.Mutex
	status int
}

func (uc *portal) Readiness() error {
	if uc.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := uc.infra.Readiness(); err != nil {
		return err
	}
	if err := uc.repos.Readiness(); err != nil {
		return err
	}
	return nil
}

func (uc *portal) Liveness() error {
	if uc.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := uc.infra.Liveness(); err != nil {
		return err
	}
	if err := uc.repos.Liveness(); err != nil {
		return err
	}
	return nil
}

func (uc *portal) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	if err := uc.infra.Connect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	uc.status = patterns.StatusConnected
	uc.logger.Info("connected")
	return nil
}

func (uc *portal) Disconnect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	uc.status = patterns.StatusDisconnected
	uc.logger.Info("disconnected")

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.infra.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *portal) Account() Account {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.account == nil {
		uc.account = &account{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.account
}

func (uc *portal) Workspace() Workspace {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.workspace == nil {
		uc.workspace = &workspace{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.workspace
}

func (uc *portal) WorkspaceCredentials() WorkspaceCredentials {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.workspaceCredentials == nil {
		uc.workspaceCredentials = &workspaceCredentials{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.workspaceCredentials
}
