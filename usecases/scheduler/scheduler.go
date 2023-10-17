package scheduler

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
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
	publisher streaming.Publisher,
	repos repos.Repositories,
) Scheduler {
	logger = logger.With("usecase", "scheduler")

	return &scheduler{
		conf:      conf,
		logger:    logger,
		infra:     infra,
		publisher: publisher,
		repos:     repos,
	}
}

type scheduler struct {
	conf      *config.Config
	logger    logging.Logger
	infra     *infrastructure.Infrastructure
	publisher streaming.Publisher
	repos     repos.Repositories

	mu      sync.RWMutex
	request *request
}

func (uc *scheduler) Readiness() error {
	if err := uc.infra.Readiness(); err != nil {
		return err
	}
	if err := uc.repos.Readiness(); err != nil {
		return err
	}
	if err := uc.publisher.Readiness(); err != nil {
		return err
	}
	return nil
}

func (uc *scheduler) Liveness() error {
	if err := uc.infra.Liveness(); err != nil {
		return err
	}
	if err := uc.repos.Liveness(); err != nil {
		return err
	}
	if err := uc.publisher.Liveness(); err != nil {
		return err
	}
	return nil
}

func (uc *scheduler) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.infra.Connect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	if err := uc.publisher.Connect(ctx); err != nil {
		return err
	}

	uc.logger.Info("connected")
	return nil
}

func (uc *scheduler) Disconnect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.logger.Info("disconnected")

	if err := uc.publisher.Disconnect(ctx); err != nil {
		return err
	}

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
			publisher: uc.publisher,
			repos:     uc.repos,
		}
	}
	return uc.request
}
