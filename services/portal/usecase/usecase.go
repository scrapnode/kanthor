package usecase

import (
	"sync"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repos"
)

type Portal interface {
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

	mu sync.Mutex
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
