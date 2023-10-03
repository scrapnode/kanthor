package scheduler

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/debugging"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/scheduler"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	subscriber streaming.Subscriber,
	metrics metric.Metrics,
	uc usecase.Scheduler,
) services.Service {
	logger = logger.With("service", "scheduler")
	return &scheduler{
		conf:       conf,
		logger:     logger,
		subscriber: subscriber,
		metrics:    metrics,
		uc:         uc,

		debugger: debugging.NewServer(),
		healthcheck: background.NewServer(
			healthcheck.DefaultConfig("scheduler"),
			logger.With("healthcheck", "background"),
		),
	}
}

type scheduler struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	metrics    metric.Metrics
	uc         usecase.Scheduler

	mu          sync.Mutex
	terminating chan bool

	debugger    debugging.Server
	healthcheck healthcheck.Server
}

func (service *scheduler) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if err := service.debugger.Start(ctx); err != nil {
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

	if err := service.healthcheck.Connect(ctx); err != nil {
		return err
	}

	service.logger.Info("started")
	return nil
}

func (service *scheduler) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.logger.Info("stopped")

	if err := service.healthcheck.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.subscriber.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.metrics.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.debugger.Stop(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *scheduler) Run(ctx context.Context) error {
	if err := service.subscriber.Sub(ctx, Consumer(service)); err != nil {
		return err
	}

	if err := service.readiness(); err != nil {
		return err
	}

	go func() {
		err := service.healthcheck.Liveness(func() error {
			service.logger.Debug("checking liveness")

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
			service.logger.Error(err)
		}
	}()

	go func() {
		if err := service.debugger.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			service.logger.Error(err)
		}
	}()

	service.logger.Info("running")
	<-service.terminating
	return nil
}

func (service *scheduler) readiness() error {
	return service.healthcheck.Readiness(func() error {
		service.logger.Debug("checking readiness")

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
