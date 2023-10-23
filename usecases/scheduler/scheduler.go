package scheduler

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

type Scheduler interface {
	patterns.Connectable
	Request() Request
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repos repos.Repositories,
) Scheduler {
	logger = logger.With("usecase", "scheduler")

	return &scheduler{
		conf:   conf,
		logger: logger,
		infra:  infra,
		repos:  repos,
	}
}

type scheduler struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories

	request *request

	mu     sync.Mutex
	status int
}

func (uc *scheduler) Readiness() error {
	if uc.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := uc.infra.Readiness(); err != nil {
		return err
	}
	if err := uc.repos.Readiness(); err != nil {
		return err
	}

	return nil
}

func (uc *scheduler) Liveness() error {
	if uc.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if err := uc.infra.Liveness(); err != nil {
		return err
	}
	if err := uc.repos.Liveness(); err != nil {
		return err
	}

	return nil
}

func (uc *scheduler) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	if err := uc.infra.Connect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	uc.status = patterns.StatusConnected
	uc.logger.Info("connected")
	return nil
}

func (uc *scheduler) Disconnect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	uc.status = patterns.StatusDisconnected
	uc.logger.Info("disconnected")

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.infra.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *scheduler) Request() Request {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.request == nil {
		uc.request = &request{
			conf:      uc.conf,
			logger:    uc.logger,
			infra:     uc.infra,
			publisher: uc.infra.Stream.Publisher("scheduler_request"),
			repos:     uc.repos,
		}
	}
	return uc.request
}
