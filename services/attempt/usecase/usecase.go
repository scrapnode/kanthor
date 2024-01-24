package usecase

import (
	"sync"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
)

type Attempt interface {
	Scanner() Scanner
	Retry() Retry
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

	scanner *scanner
	retry   *retry

	mu sync.Mutex
}

func (uc *attempt) Scanner() Scanner {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.scanner == nil {
		uc.scanner = &scanner{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			publisher:    uc.infra.Stream.Publisher("attempt.scanner"),
			repositories: uc.repositories,
		}
	}
	return uc.scanner
}

func (uc *attempt) Retry() Retry {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.retry == nil {
		uc.retry = &retry{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			publisher:    uc.infra.Stream.Publisher("attempt.retry"),
			repositories: uc.repositories,
		}
	}
	return uc.retry
}
