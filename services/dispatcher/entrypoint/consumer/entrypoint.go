package consumer

import (
	"context"
	"errors"
	"sync"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/constants"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/project"
	"github.com/scrapnode/kanthor/services/dispatcher/config"
	"github.com/scrapnode/kanthor/services/dispatcher/usecase"
	"github.com/scrapnode/kanthor/telemetry"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	uc usecase.Dispatcher,
) patterns.Runnable {
	logger = logger.With("service", "dispatcher", "entrypoint", "consumer")
	return &dispatcher{
		conf:       conf,
		logger:     logger,
		subscriber: infra.Stream.Subscriber("dispatcher_consumer"),
		infra:      infra,
		uc:         uc,

		healthcheck: background.NewServer(
			healthcheck.DefaultConfig("dispatcher.consumer"),
			logger.With("healthcheck", "background"),
		),
	}
}

type dispatcher struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	infra      *infrastructure.Infrastructure
	uc         usecase.Dispatcher

	healthcheck healthcheck.Server

	mu     sync.Mutex
	status int
}

func (service *dispatcher) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	if err := service.conf.Validate(); err != nil {
		return err
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	if err := service.subscriber.Connect(ctx); err != nil {
		return err
	}

	if err := service.healthcheck.Connect(ctx); err != nil {
		return err
	}

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *dispatcher) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status != patterns.StatusStarted {
		return ErrNotStarted
	}
	service.status = patterns.StatusStopped
	service.logger.Info("stopped")

	var returning error
	if err := service.healthcheck.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.subscriber.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.infra.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	return returning
}

func (service *dispatcher) Run(ctx context.Context) error {
	tracectx := context.WithValue(ctx, telemetry.CtxTracer, telemetry.Tracer(project.Name("dispatcher_consumer")))
	topic := constants.TopicRequest
	if err := service.subscriber.Sub(tracectx, topic, Handler(service)); err != nil {
		return err
	}

	if err := service.readiness(); err != nil {
		return err
	}

	go func() {
		err := service.healthcheck.Liveness(func() error {
			if err := service.subscriber.Liveness(); err != nil {
				return err
			}

			if err := service.infra.Liveness(); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			service.logger.Error(err)
		}
	}()

	service.logger.Infow("running", "topic", topic)
	<-ctx.Done()
	return nil
}

func (service *dispatcher) readiness() error {
	return service.healthcheck.Readiness(func() error {
		if err := service.subscriber.Readiness(); err != nil {
			return err
		}

		if err := service.infra.Readiness(); err != nil {
			return err
		}

		return nil
	})
}
