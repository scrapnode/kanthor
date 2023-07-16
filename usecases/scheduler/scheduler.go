package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
	"sync"
)

type Scheduler interface {
	patterns.Connectable
	Request() Request
}

func New(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	cache cache.Cache,
	meter metric.Meter,
	repos repos.Repositories,
) Scheduler {
	return &scheduler{
		conf:      conf,
		logger:    logger,
		timer:     timer,
		publisher: publisher,
		cache:     cache,
		meter:     meter,
		repos:     repos,
	}
}

type scheduler struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	cache     cache.Cache
	meter     metric.Meter
	repos     repos.Repositories

	mu      sync.RWMutex
	request *request
}

func (usecase *scheduler) Connect(ctx context.Context) error {
	if err := usecase.repos.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.publisher.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Connect(ctx); err != nil {
		return err
	}

	usecase.logger.Info("connected")
	return nil
}

func (usecase *scheduler) Disconnect(ctx context.Context) error {
	usecase.logger.Info("disconnected")

	if err := usecase.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.publisher.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (usecase *scheduler) Request() Request {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.request == nil {
		usecase.request = &request{
			conf:      usecase.conf,
			logger:    usecase.logger,
			timer:     usecase.timer,
			publisher: usecase.publisher,
			repos:     usecase.repos,
			cache:     usecase.cache,
			meter:     usecase.meter,
		}
	}

	return usecase.request
}
