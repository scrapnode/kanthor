package scheduler

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/signature"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

type Scheduler interface {
	patterns.Connectable
	Request() Request
}

func New(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	signature signature.Signature,
	publisher streaming.Publisher,
	cache cache.Cache,
	metrics metric.Metrics,
	repos repos.Repositories,
) Scheduler {
	return &scheduler{
		conf:      conf,
		logger:    logger,
		timer:     timer,
		signature: signature,
		publisher: publisher,
		cache:     cache,
		metrics:   metrics,
		repos:     repos,
	}
}

type scheduler struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	signature signature.Signature
	publisher streaming.Publisher
	cache     cache.Cache
	metrics   metric.Metrics
	repos     repos.Repositories

	mu      sync.RWMutex
	request *request
}

func (uc *scheduler) Readiness() error {
	if err := uc.cache.Readiness(); err != nil {
		return err
	}
	if err := uc.repos.Readiness(); err != nil {
		return err
	}
	if err := uc.publisher.Readiness(); err != nil {
		return err
	}
	return nil
}

func (uc *scheduler) Liveness() error {
	if err := uc.cache.Liveness(); err != nil {
		return err
	}
	if err := uc.repos.Liveness(); err != nil {
		return err
	}
	if err := uc.publisher.Liveness(); err != nil {
		return err
	}
	return nil
}

func (uc *scheduler) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.cache.Connect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	if err := uc.publisher.Connect(ctx); err != nil {
		return err
	}

	uc.logger.Info("connected")
	return nil
}

func (uc *scheduler) Disconnect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.logger.Info("disconnected")

	if err := uc.publisher.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *scheduler) Request() Request {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.request == nil {
		uc.request = &request{
			conf:      uc.conf,
			logger:    uc.logger,
			timer:     uc.timer,
			signature: uc.signature,
			publisher: uc.publisher,
			cache:     uc.cache,
			metrics:   uc.metrics,
			repos:     uc.repos,
		}
	}
	return uc.request
}
