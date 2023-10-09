package attempt

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/attempt/repos"
)

type Attempt interface {
	patterns.Connectable
	Application() Application
	Trigger() Trigger
}

func New(
	conf *config.Config,
	logger logging.Logger,
	cache cache.Cache,
	timer timer.Timer,
	publisher streaming.Publisher,
	metrics metric.Metrics,
	repos repos.Repositories,
) Attempt {
	logger = logger.With("usecase", "attempt")

	return &attempt{
		conf:      conf,
		logger:    logger,
		cache:     cache,
		timer:     timer,
		publisher: publisher,
		metrics:   metrics,
		repos:     repos,
	}
}

type attempt struct {
	conf      *config.Config
	logger    logging.Logger
	cache     cache.Cache
	timer     timer.Timer
	publisher streaming.Publisher
	metrics   metric.Metrics
	repos     repos.Repositories

	mu          sync.RWMutex
	application *application
	trigger     *trigger
}

func (uc *attempt) Readiness() error {
	if err := uc.repos.Readiness(); err != nil {
		return err
	}
	return nil
}

func (uc *attempt) Liveness() error {
	if err := uc.repos.Liveness(); err != nil {
		return err
	}
	return nil
}

func (uc *attempt) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	uc.logger.Info("connected")
	return nil
}

func (uc *attempt) Disconnect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.logger.Info("disconnected")

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *attempt) Application() Application {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.application == nil {
		uc.application = &application{
			conf:      uc.conf,
			logger:    uc.logger,
			cache:     uc.cache,
			timer:     uc.timer,
			publisher: uc.publisher,
			metrics:   uc.metrics,
			repos:     uc.repos,
		}
	}
	return uc.application
}

func (uc *attempt) Trigger() Trigger {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.trigger == nil {
		uc.trigger = &trigger{
			conf:      uc.conf,
			logger:    uc.logger,
			cache:     uc.cache,
			timer:     uc.timer,
			publisher: uc.publisher,
			metrics:   uc.metrics,
			repos:     uc.repos,
		}
	}
	return uc.trigger
}
