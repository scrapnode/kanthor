package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/usecases"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	subscriber streaming.Subscriber,
	usecase usecases.Scheduler,
) services.Service {
	logger.With("service", "scheduler")
	return &scheduler{conf: conf, logger: logger, subscriber: subscriber, usecase: usecase}
}

type scheduler struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	usecase    usecases.Scheduler
}

func (service *scheduler) Start(ctx context.Context) error {
	if err := service.usecase.Connect(ctx); err != nil {
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

	if err := service.usecase.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *scheduler) Run(ctx context.Context) error {
	logger := service.logger.With("fn", "consumer")
	return service.subscriber.Sub(ctx, Consumer(logger, service.usecase))
}
