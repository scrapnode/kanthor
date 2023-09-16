package dispatcher

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/debugging"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	validator validator.Validator,
	subscriber streaming.Subscriber,
	metrics metric.Metrics,
	uc usecase.Dispatcher,
) services.Service {
	logger = logger.With("service", "dispatcher")
	return &dispatcher{
		conf:       conf,
		logger:     logger,
		validator:  validator,
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
	validator  validator.Validator
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

	if err := service.debugger.Stop(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *dispatcher) Run(ctx context.Context) error {
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
		if err := service.debugger.Run(ctx); err != nil {
			service.logger.Error(err)
		}
	}()

	service.logger.Info("running")
	return service.subscriber.Sub(ctx, Consumer(service))
}

func (service *dispatcher) readiness() error {
	return service.healthcheck.Readiness(func() error {
		// @TODO: add starting up checking here
		return nil
	})
}
