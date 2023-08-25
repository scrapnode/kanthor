package dispatcher

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics/exporter"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
	"net/http"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	subscriber streaming.Subscriber,
	metrics metrics.Metrics,
	uc usecase.Dispatcher,
) services.Service {
	logger = logger.With("service", "dispatcher")
	return &dispatcher{
		conf:       conf,
		logger:     logger,
		subscriber: subscriber,
		metrics:    metrics,
		exporter:   metrics.Exporter(),
		uc:         uc,
	}
}

type dispatcher struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	metrics    metrics.Metrics
	exporter   exporter.Exporter
	uc         usecase.Dispatcher
}

func (service *dispatcher) Start(ctx context.Context) error {
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

func (service *dispatcher) Stop(ctx context.Context) error {
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

func (service *dispatcher) Run(ctx context.Context) error {
	go func() {
		if err := service.exporter.Run(ctx); err != nil && err != http.ErrServerClosed {
			service.logger.Error(err)
		}
	}()

	service.logger.Info("running")
	return service.subscriber.Sub(ctx, Consumer(service))
}
