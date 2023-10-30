package planner

import (
	"context"
	"errors"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/project"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	logger = logger.With("service", "attempt.endeavor.planner")
	return &planner{
		conf:   conf,
		logger: logger,
		infra:  infra,
		ds:     ds,
		uc:     uc,

		cron: cron.New(),
		healthcheck: background.NewServer(
			healthcheck.DefaultConfig("attempt.endeavor.planner"),
			logger.With("healthcheck", "background"),
		),
	}
}

type planner struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	ds     datastore.Datastore
	uc     usecase.Attempt

	cron        *cron.Cron
	healthcheck healthcheck.Server

	mu     sync.Mutex
	status int
}

func (service *planner) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	if err := service.ds.Connect(ctx); err != nil {
		return err
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	if err := service.healthcheck.Connect(ctx); err != nil {
		return err
	}

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *planner) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status != patterns.StatusStarted {
		return ErrNotStarted
	}
	service.status = patterns.StatusStopped
	service.logger.Info("stopped")

	cronctx := service.cron.Stop()
	// wait for all processing jobs completed
	<-cronctx.Done()

	var returning error
	if err := service.healthcheck.Disconnect(ctx); err != nil {
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

func (service *planner) Run(ctx context.Context) error {
	id, err := service.cron.AddFunc(service.conf.Trigger.Planner.Schedule, RegisterCron(service))
	if err != nil {
		return err
	}

	// on dev environment, should run the job immediately after we starting the service
	if project.IsDev() {
		service.logger.Debug("starting immediately because of development env")
		service.cron.Entry(id).Job.Run()
	}

	if err := service.readiness(); err != nil {
		return err
	}

	go func() {
		err := service.healthcheck.Liveness(func() error {
			service.logger.Debug("checking liveness")

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

	service.logger.Info("running", "schedule", service.conf.Trigger.Planner.Schedule)
	service.cron.Run()
	return nil
}

func (service *planner) readiness() error {
	return service.healthcheck.Readiness(func() error {
		service.logger.Debug("checking readiness")

		if err := service.infra.Readiness(); err != nil {
			return err
		}

		if err := service.ds.Readiness(); err != nil {
			return err
		}

		return nil
	})
}
