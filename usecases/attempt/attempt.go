package attempt

import (
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/usecases/attempt/repos"
)

type Attempt interface {
	Trigger() Trigger
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repos repos.Repositories,
) Attempt {
	logger = logger.With("usecase", "attempt")

	return &attempt{
		conf:   conf,
		logger: logger,
		infra:  infra,
		repos:  repos,
	}
}

type attempt struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories

	trigger *trigger

	mu sync.Mutex
}

func (uc *attempt) Trigger() Trigger {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.trigger == nil {
		uc.trigger = &trigger{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.trigger
}
