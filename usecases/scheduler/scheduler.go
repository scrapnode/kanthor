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
	uc := &scheduler{
		conf:      conf,
		logger:    logger,
		timer:     timer,
		publisher: publisher,
		cache:     cache,
		meter:     meter,
		repos:     repos,
	}

	uc.request = &request{
		conf:      uc.conf,
		logger:    uc.logger,
		timer:     uc.timer,
		publisher: uc.publisher,
		repos:     uc.repos,
		cache:     uc.cache,
		meter:     uc.meter,
	}
	
	return uc
}

type scheduler struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	cache     cache.Cache
	meter     metric.Meter
	repos     repos.Repositories

	request *request
}

func (uc *scheduler) Connect(ctx context.Context) error {
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
	return uc.request
}
