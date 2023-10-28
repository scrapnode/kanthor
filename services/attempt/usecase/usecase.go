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

	trigger *trigger

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
