package dispatcher

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
)

type Dispatcher interface {
	patterns.Connectable
	Forwarder() Forwarder
}

func New(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	dispatch sender.Send,
	cache cache.Cache,
	cb circuitbreaker.CircuitBreaker,
	metrics metric.Metrics,
) Dispatcher {
	return &dispatcher{
		conf:      conf,
		logger:    logger,
		timer:     timer,
		publisher: publisher,
		dispatch:  dispatch,
		cache:     cache,
		cb:        cb,
		metrics:   metrics,
	}
}

type dispatcher struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	dispatch  sender.Send
	cache     cache.Cache
	cb        circuitbreaker.CircuitBreaker
	metrics   metric.Metrics

	mu        sync.RWMutex
	forwarder *forwarder
}

func (uc *dispatcher) Readiness() error {
	if err := uc.cache.Readiness(); err != nil {
		return err
	}
	if err := uc.publisher.Readiness(); err != nil {
		return err
	}
	return nil
}

func (uc *dispatcher) Liveness() error {
	if err := uc.cache.Liveness(); err != nil {
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

	if err := uc.cache.Connect(ctx); err != nil {
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

	if err := uc.cache.Disconnect(ctx); err != nil {
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
			timer:     uc.timer,
			publisher: uc.publisher,
			dispatch:  uc.dispatch,
			cache:     uc.cache,
			cb:        uc.cb,
			metrics:   uc.metrics,
		}
	}
	return uc.forwarder
}
