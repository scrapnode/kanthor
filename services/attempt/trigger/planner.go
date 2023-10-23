package trigger

import (
	"context"
	"errors"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/attempt"
)

func NewPlanner(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	uc usecase.Attempt,
) services.Service {
	logger = logger.With("service", "attempt.trigger.planner")
	return &planner{
		conf:   conf,
		logger: logger,
		infra:  infra,
		uc:     uc,

		cron: cron.New(),
		healthcheck: background.NewServer(
			healthcheck.DefaultConfig("attempt.trigger.planner"),
			logger.With("healthcheck", "background"),
		),
	}
}

type planner struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	uc     usecase.Attempt

	cron        *cron.Cron
	healthcheck healthcheck.Server

	mu      sync.Mutex
	status  int
	stopped chan bool
}

func (service *planner) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrPlannerAlreadyStarted
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	if err := service.uc.Connect(ctx); err != nil {
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
		return ErrPlannerNotStarted
	}
	service.status = patterns.StatusStopped
	service.logger.Info("stopped")

	cronctx := service.cron.Stop()
	// wait for all processing jobs completed
	<-cronctx.Done()

	var returning error
	if err := service.healthcheck.Disconnect(ctx); err != nil {
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

	service.stopped <- true
	return returning
}

func (service *planner) Run(ctx context.Context) error {
	id, err := service.cron.AddFunc(service.conf.Attempt.Trigger.Planner.Schedule, RegisterCron(service))
	if err != nil {
		return err
	}
	service.cron.Run()
	// on dev environment, should run the job immediately after we starting the service
	if service.conf.Development {
		service.cron.Entry(id).Job.Run()
	}

	if err := service.readiness(); err != nil {
		return err
	}

	go func() {
		err := service.healthcheck.Liveness(func() error {
			service.logger.Debug("checking liveness")

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

	service.logger.Info("running")
	<-service.stopped
	return nil
}

func (service *planner) readiness() error {
	return service.healthcheck.Readiness(func() error {
		service.logger.Debug("checking readiness")

		if err := service.uc.Readiness(); err != nil {
			return err
		}

		if err := service.infra.Readiness(); err != nil {
			return err
		}
		return nil
	})
}
