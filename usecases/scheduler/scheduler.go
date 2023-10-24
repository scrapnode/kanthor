package scheduler

import (
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

type Scheduler interface {
	Request() Request
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repos repos.Repositories,
) Scheduler {
	logger = logger.With("usecase", "scheduler")

	return &scheduler{
		conf:   conf,
		logger: logger,
		infra:  infra,
		repos:  repos,
	}
}

type scheduler struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories

	request *request

	mu sync.Mutex
}

func (uc *scheduler) Request() Request {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.request == nil {
		uc.request = &request{
			conf:      uc.conf,
			logger:    uc.logger,
			infra:     uc.infra,
			publisher: uc.infra.Stream.Publisher("scheduler_request"),
			repos:     uc.repos,
		}
	}
	return uc.request
}
