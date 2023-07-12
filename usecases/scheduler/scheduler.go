package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
	"sync"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	repos repos.Repositories,
	cache cache.Cache,
	meter metric.Meter,
) Scheduler {
	return &scheduler{
		conf:      conf,
		logger:    logger,
		timer:     timer,
		publisher: publisher,
		repos:     repos,
		cache:     cache,
		meter:     meter,
	}
}

type scheduler struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	repos     repos.Repositories
	cache     cache.Cache
	meter     metric.Meter

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
