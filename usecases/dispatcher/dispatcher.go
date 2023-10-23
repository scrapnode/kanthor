package dispatcher

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Dispatcher interface {
	patterns.Connectable
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

	mu     sync.Mutex
	status int
}

func (uc *dispatcher) Readiness() error {
	if uc.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := uc.infra.Readiness(); err != nil {
		return err
	}

	return nil
}

func (uc *dispatcher) Liveness() error {
	if uc.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := uc.infra.Liveness(); err != nil {
		return err
	}

	return nil
}

func (uc *dispatcher) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	if err := uc.infra.Connect(ctx); err != nil {
		return err
	}

	uc.status = patterns.StatusConnected
	uc.logger.Info("connected")
	return nil
}

func (uc *dispatcher) Disconnect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	uc.status = patterns.StatusDisconnected
	uc.logger.Info("disconnected")

	if err := uc.infra.Disconnect(ctx); err != nil {
		return err
	}

	return nil
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
