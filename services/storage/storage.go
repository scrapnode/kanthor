package storage

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/storage"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	subscriber streaming.Subscriber,
	metrics metric.Metrics,
	uc usecase.Storage,
) services.Service {
	logger = logger.With("service", "storage")
	return &storage{
		conf:       conf,
		logger:     logger,
		subscriber: subscriber,
		metrics:    metrics,
		uc:         uc,

		healthcheck: background.NewServer(
			healthcheck.DefaultConfig("kanthor.storage"),
			logger.With("healthcheck", "background"),
		),
	}
}

type storage struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	metrics    metric.Metrics
	uc         usecase.Storage

	healthcheck healthcheck.Server
}

func (service *storage) Start(ctx context.Context) error {
	if err := service.metrics.Connect(ctx); err != nil {
		return err
	}

	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	if err := service.subscriber.Connect(ctx); err != nil {
		return err
	}

	service.logger.Info("started")
	return nil
}

func (service *storage) Stop(ctx context.Context) error {
	service.logger.Info("stopped")

	if err := service.subscriber.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.metrics.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *storage) Run(ctx context.Context) error {
	if err := service.readiness(); err != nil {
		return err
	}

	go func() {
		err := service.healthcheck.Liveness(func() error {
			return nil
		})
		if err != nil {
			service.logger.Error(err)
		}
	}()

	service.logger.Info("running")
	return service.subscriber.Sub(ctx, Consumer(service))
}

func (service *storage) readiness() error {
	return service.healthcheck.Readiness(func() error {
		// @TODO: add starting up checking here
		return nil
	})
}
