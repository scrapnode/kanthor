package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/debugging"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	subscriber streaming.Subscriber,
	metrics metric.Metrics,
	uc usecase.Dispatcher,
) services.Service {
	logger = logger.With("service", "dispatcher")
	return &dispatcher{
		conf:       conf,
		logger:     logger,
		subscriber: subscriber,
		metrics:    metrics,
		uc:         uc,

		debugger: debugging.NewServer(),
		healthcheck: background.NewServer(
			healthcheck.DefaultConfig("dispatcher"),
			logger.With("healthcheck", "background"),
		),
	}
}

type dispatcher struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	metrics    metric.Metrics
	uc         usecase.Dispatcher

	debugger    debugging.Server
	healthcheck healthcheck.Server
}

func (service *dispatcher) Start(ctx context.Context) error {
	if err := service.debugger.Start(ctx); err != nil {
		return err
	}

	if err := service.healthcheck.Connect(ctx); err != nil {
		return err
	}

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

func (service *dispatcher) Stop(ctx context.Context) error {
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

	if err := service.healthcheck.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.debugger.Stop(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *dispatcher) Run(ctx context.Context) error {
	if err := service.readiness(); err != nil {
		return fmt.Errorf("HEALTHCHECK.READINESS: %v", err)
	}

	go func() {
		err := service.healthcheck.Liveness(func() error {
			if err := service.subscriber.Liveness(); err != nil {
				return err
			}

			if err := service.uc.Liveness(); err != nil {
				return err
			}

			if err := service.metrics.Liveness(); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			service.logger.Errorf("HEALTHCHECK.LIVENESS: %v", err)
		}
	}()

	go func() {
		if err := service.debugger.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			service.logger.Error(err)
		}
	}()

	service.logger.Info("running")
	return service.subscriber.Sub(ctx, Consumer(service))
}

func (service *dispatcher) readiness() error {
	return service.healthcheck.Readiness(func() error {
		if err := service.subscriber.Readiness(); err != nil {
			return err
		}

		if err := service.uc.Readiness(); err != nil {
			return err
		}

		if err := service.metrics.Readiness(); err != nil {
			return err
		}

		return nil
	})
}
