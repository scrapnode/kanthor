package dispatcher

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"sync"
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
	meter metric.Meter,
) Dispatcher {
	return &dispatcher{
		conf:      conf,
		logger:    logger,
		timer:     timer,
		publisher: publisher,
		dispatch:  dispatch,
		cache:     cache,
		cb:        cb,
		meter:     meter,
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
	meter     metric.Meter

	mu        sync.RWMutex
	forwarder *forwarder
}

func (usecase *dispatcher) Connect(ctx context.Context) error {
	if err := usecase.publisher.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Connect(ctx); err != nil {
		return err
	}

	usecase.logger.Info("connected")
	return nil
}

func (usecase *dispatcher) Disconnect(ctx context.Context) error {
	usecase.logger.Info("disconnected")

	if err := usecase.publisher.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (usecase *dispatcher) Forwarder() Forwarder {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.forwarder == nil {
		usecase.forwarder = &forwarder{
			conf:      usecase.conf,
			logger:    usecase.logger,
			timer:     usecase.timer,
			publisher: usecase.publisher,
			cache:     usecase.cache,
			meter:     usecase.meter,
		}
	}

	return usecase.forwarder
}
