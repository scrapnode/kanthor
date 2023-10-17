package dispatcher

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
)

type Dispatcher interface {
	patterns.Connectable
	Forwarder() Forwarder
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	publisher streaming.Publisher,
	dispatch sender.Send,
) Dispatcher {
	logger = logger.With("usecase", "dispatcher")

	return &dispatcher{
		conf:      conf,
		logger:    logger,
		infra:     infra,
		publisher: publisher,
		dispatch:  dispatch,
	}
}

type dispatcher struct {
	conf      *config.Config
	logger    logging.Logger
	infra     *infrastructure.Infrastructure
	dispatch  sender.Send
	publisher streaming.Publisher

	mu        sync.RWMutex
	forwarder *forwarder
}

func (uc *dispatcher) Readiness() error {
	if err := uc.infra.Readiness(); err != nil {
		return err
	}
	if err := uc.publisher.Readiness(); err != nil {
		return err
	}
	return nil
}

func (uc *dispatcher) Liveness() error {
	if err := uc.infra.Liveness(); err != nil {
		return err
	}
	if err := uc.publisher.Liveness(); err != nil {
		return err
	}
	return nil
}

func (uc *dispatcher) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.infra.Connect(ctx); err != nil {
		return err
	}

	if err := uc.publisher.Connect(ctx); err != nil {
		return err
	}

	uc.logger.Info("connected")
	return nil
}

func (uc *dispatcher) Disconnect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.logger.Info("disconnected")

	if err := uc.publisher.Disconnect(ctx); err != nil {
		return err
	}

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
			publisher: uc.publisher,
			dispatch:  uc.dispatch,
		}
	}
	return uc.forwarder
}
