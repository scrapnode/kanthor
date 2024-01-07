package consumer

import (
	"context"
	"errors"
	"sync"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/constants"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services/storage/config"
	"github.com/scrapnode/kanthor/services/storage/usecase"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	ds datastore.Datastore,
	uc usecase.Storage,
) patterns.Runnable {
	logger = logger.With("service", "storage", "entrypoint", "consumer")
	return &storage{
		conf:       conf,
		logger:     logger,
		subscriber: infra.Stream.Subscriber("storage"),
		infra:      infra,
		ds:         ds,
		uc:         uc,

		healthcheck: background.NewServer(
			healthcheck.DefaultConfig("storage"),
			logger.With("healthcheck", "background"),
		),
	}
}

type storage struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	infra      *infrastructure.Infrastructure
	ds         datastore.Datastore
	uc         usecase.Storage

	healthcheck healthcheck.Server

	mu     sync.Mutex
	status int
}

func (service *storage) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	if err := service.conf.Validate(); err != nil {
		return err
	}

	if err := service.ds.Connect(ctx); err != nil {
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

func (service *storage) Stop(ctx context.Context) error {
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

	if err := service.ds.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	return returning
}

func (service *storage) Run(ctx context.Context) error {
	topic := constants.TopicPublic
	if service.conf.Topic != "" {
		topic = service.conf.Topic
	}

	if err := service.subscriber.Sub(ctx, topic, Handler(service)); err != nil {
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

			if err := service.infra.Liveness(); err != nil {
				return err
			}

			if err := service.ds.Liveness(); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			service.logger.Error(err)
		}
	}()

	service.logger.Infow("running", "topic", topic)
	forever := make(chan bool, 1)
	select {
	case <-forever:
		return nil
	case <-ctx.Done():
		return nil
	}
}

func (service *storage) readiness() error {
	return service.healthcheck.Readiness(func() error {
		service.logger.Debug("checking readiness")

		if err := service.subscriber.Readiness(); err != nil {
			return err
		}

		if err := service.infra.Readiness(); err != nil {
			return err
		}

		if err := service.ds.Readiness(); err != nil {
			return err
		}

		return nil
	})
}
