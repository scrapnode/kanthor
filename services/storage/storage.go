package storage

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/debugging"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/storage"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	subscriber streaming.Subscriber,
	uc usecase.Storage,
) services.Service {
	logger = logger.With("service", "storage")
	return &storage{
		consumer:   "storage",
		conf:       conf,
		logger:     logger,
		subscriber: subscriber,
		infra:      infra,
		uc:         uc,

		debugger: debugging.NewServer(),
		healthcheck: background.NewServer(
			healthcheck.DefaultConfig("storage"),
			logger.With("healthcheck", "background"),
		),
	}
}

type storage struct {
	consumer string

	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	infra      *infrastructure.Infrastructure
	uc         usecase.Storage

	mu          sync.Mutex
	terminating chan bool

	debugger    debugging.Server
	healthcheck healthcheck.Server
}

func (service *storage) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if err := service.debugger.Start(ctx); err != nil {
		return err
	}
	if err := service.infra.Metric.Connect(ctx); err != nil {
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

func (service *storage) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.logger.Info("stopped")
	var returning error

	if err := service.healthcheck.Disconnect(ctx); err != nil {
		service.logger.Error(err)
		returning = errors.Join(returning, err)
	}

	if err := service.subscriber.Disconnect(ctx); err != nil {
		service.logger.Error(err)
		returning = errors.Join(returning, err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
		returning = errors.Join(returning, err)
	}

	if err := service.infra.Disconnect(ctx); err != nil {
		service.logger.Error(err)
		returning = errors.Join(returning, err)
	}

	if err := service.debugger.Stop(ctx); err != nil {
		service.logger.Error(err)
		returning = errors.Join(returning, err)
	}

	return returning
}

func (service *storage) Run(ctx context.Context) error {
	if err := service.subscriber.Sub(ctx, service.consumer, constants.TopicPublic, NewConsumer(service)); err != nil {
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

			if err := service.infra.Liveness(); err != nil {
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

func (service *storage) readiness() error {
	return service.healthcheck.Readiness(func() error {
		service.logger.Debug("checking readiness")

		if err := service.subscriber.Readiness(); err != nil {
			return err
		}

		if err := service.uc.Readiness(); err != nil {
			return err
		}

		if err := service.infra.Readiness(); err != nil {
			return err
		}
		return nil
	})
}
