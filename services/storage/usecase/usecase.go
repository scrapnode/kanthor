package usecase

import (
	"sync"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/storage/config"
	"github.com/scrapnode/kanthor/services/storage/repos"
)

type Storage interface {
	Warehouse() Warehouse
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repos repos.Repositories,
) Storage {
	logger = logger.With("usecase", "storage")

	return &storage{
		conf:   conf,
		logger: logger,
		infra:  infra,
		repos:  repos,
	}
}

type storage struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories

	warehose *warehose

	mu sync.Mutex
}

func (uc *storage) Warehouse() Warehouse {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.warehose == nil {
		uc.warehose = &warehose{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.warehose
}
