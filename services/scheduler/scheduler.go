package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/scheduler"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	subscriber streaming.Subscriber,
	uc usecase.Scheduler,
	meter metric.Meter,
) services.Service {
	logger = logger.With("service", "scheduler")
	return &scheduler{conf: conf, logger: logger, subscriber: subscriber, uc: uc, meter: meter}
}

type scheduler struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	uc         usecase.Scheduler
	meter      metric.Meter
}

func (service *scheduler) Start(ctx context.Context) error {
	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	if err := service.subscriber.Connect(ctx); err != nil {
		return err
	}

	service.logger.Info("started")
	return nil
}

func (service *scheduler) Stop(ctx context.Context) error {
	service.logger.Info("stopped")

	if err := service.subscriber.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *scheduler) Run(ctx context.Context) error {
	return service.subscriber.Sub(ctx, Consumer(service))
}
