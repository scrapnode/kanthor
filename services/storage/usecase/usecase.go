package usecase

import (
	"sync"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/storage/config"
	"github.com/scrapnode/kanthor/services/storage/repositories"
)

type Storage interface {
	Warehouse() Warehouse
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repositories repositories.Repositories,
) Storage {
	logger = logger.With("usecase", "storage")

	return &storage{
		conf:         conf,
		logger:       logger,
		infra:        infra,
		repositories: repositories,
	}
}

type storage struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories

	warehose *warehose

	mu sync.Mutex
}

func (uc *storage) Warehouse() Warehouse {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.warehose == nil {
		uc.warehose = &warehose{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			repositories: uc.repositories,
		}
	}
	return uc.warehose
}
