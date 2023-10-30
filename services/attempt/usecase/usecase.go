package usecase

import (
	"sync"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
)

type Attempt interface {
	Trigger() Trigger
	Endeavor() Endeavor
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repositories repositories.Repositories,
) Attempt {
	logger = logger.With("usecase", "attempt")

	return &attempt{
		conf:         conf,
		logger:       logger,
		infra:        infra,
		repositories: repositories,
	}
}

type attempt struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories

	trigger  *trigger
	endeavor *endeavor

	mu sync.Mutex
}

func (uc *attempt) Trigger() Trigger {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.trigger == nil {
		uc.trigger = &trigger{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			repositories: uc.repositories,
		}
	}
	return uc.trigger
}

func (uc *attempt) Endeavor() Endeavor {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.endeavor == nil {
		uc.endeavor = &endeavor{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			repositories: uc.repositories,
			publisher:    uc.infra.Stream.Publisher("attempt_endeavor"),
		}
	}
	return uc.endeavor
}
