package usecase

import (
	"sync"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/scheduler/config"
	"github.com/scrapnode/kanthor/services/scheduler/repositories"
)

type Scheduler interface {
	Request() Request
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repositories repositories.Repositories,
) Scheduler {
	logger = logger.With("usecase", "scheduler")

	return &scheduler{
		conf:         conf,
		logger:       logger,
		infra:        infra,
		repositories: repositories,
	}
}

type scheduler struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories

	request *request

	mu sync.Mutex
}

func (uc *scheduler) Request() Request {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.request == nil {
		uc.request = &request{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			publisher:    uc.infra.Stream.Publisher("scheduler_request"),
			repositories: uc.repositories,
		}
	}
	return uc.request
}
