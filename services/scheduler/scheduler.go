package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics/exporter"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/scheduler"
	"net/http"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	subscriber streaming.Subscriber,
	metrics metrics.Metrics,
	uc usecase.Scheduler,
) services.Service {
	logger = logger.With("service", "scheduler")
	return &scheduler{
		conf:       conf,
		logger:     logger,
		subscriber: subscriber,
		metrics:    metrics,
		exporter:   metrics.Exporter(),
		uc:         uc,

		healthcheck: background.NewServer(healthcheck.DefaultConfig("kanthor.scheduler")),
	}
}

type scheduler struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	metrics    metrics.Metrics
	exporter   exporter.Exporter
	uc         usecase.Scheduler

	healthcheck healthcheck.Server
}

func (service *scheduler) Start(ctx context.Context) error {
	if err := service.exporter.Start(ctx); err != nil {
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

func (service *scheduler) Stop(ctx context.Context) error {
	service.logger.Info("stopped")

	if err := service.subscriber.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.exporter.Stop(ctx); err != nil {
		return err
	}

	return nil
}

func (service *scheduler) Run(ctx context.Context) error {
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

	go func() {
		if err := service.exporter.Run(ctx); err != nil && err != http.ErrServerClosed {
			service.logger.Error(err)
		}
	}()

	service.logger.Info("running")
	return service.subscriber.Sub(ctx, Consumer(service))
}

func (service *scheduler) readiness() error {
	return service.healthcheck.Readiness(func() error {
		// @TODO: add starting up checking here
		return nil
	})
}
