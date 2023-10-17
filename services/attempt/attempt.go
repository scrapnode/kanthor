package attempt

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/debugging"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/attempt"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	uc usecase.Attempt,
) services.Service {
	logger = logger.With("service", "attempt")
	return &attempt{
		conf:   conf,
		logger: logger,
		infra:  infra,
		uc:     uc,

		cron:     cron.New(),
		debugger: debugging.NewServer(),
		healthcheck: background.NewServer(
			healthcheck.DefaultConfig("attempt"),
			logger.With("healthcheck", "background"),
		),
	}
}

type attempt struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	uc     usecase.Attempt

	mu          sync.Mutex
	terminating chan bool

	cron        *cron.Cron
	debugger    debugging.Server
	healthcheck healthcheck.Server
}

func (service *attempt) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if err := service.debugger.Start(ctx); err != nil {
		return err
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

	// robfig/cron provide a function cron.Start() but it's just a shortcut of cron.Run()
	service.cron.AddFunc("", RegisterTriggerCron(service))
	service.logger.Info("started")
	return nil
}

func (service *attempt) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.logger.Info("stopped")

	cronctx := service.cron.Stop()
	// wait for all processing jobs completed
	<-cronctx.Done()

	if err := service.healthcheck.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.infra.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.debugger.Stop(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *attempt) Run(ctx context.Context) error {
	service.cron.Run()

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

	go func() {
		if err := service.debugger.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			service.logger.Error(err)
		}
	}()

	service.logger.Info("running")
	<-service.terminating
	return nil
}

func (service *attempt) readiness() error {
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
