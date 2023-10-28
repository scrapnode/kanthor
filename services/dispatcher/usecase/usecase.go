package usecase

import (
	"sync"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/dispatcher/config"
)

type Dispatcher interface {
	Forwarder() Forwarder
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
) Dispatcher {
	logger = logger.With("usecase", "dispatcher")

	return &dispatcher{
		conf:   conf,
		logger: logger,
		infra:  infra,
	}
}

type dispatcher struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure

	forwarder *forwarder

	mu sync.Mutex
}

func (uc *dispatcher) Forwarder() Forwarder {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.forwarder == nil {
		uc.forwarder = &forwarder{
			conf:      uc.conf,
			logger:    uc.logger,
			infra:     uc.infra,
			publisher: uc.infra.Stream.Publisher("dispatcher_fowarder"),
		}
	}
	return uc.forwarder
}
