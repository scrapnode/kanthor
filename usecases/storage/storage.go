package storage

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

type Storage interface {
	patterns.Connectable
	Warehouse() Warehouse
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repos repos.Repositories,
) Storage {
	logger = logger.With("usecase", "storage")

	return &storage{
		conf:   conf,
		logger: logger,
		infra:  infra,
		repos:  repos,
	}
}

type storage struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories

	mu       sync.RWMutex
	warehose *warehose
}

func (uc *storage) Readiness() error {
	if err := uc.repos.Readiness(); err != nil {
		return err
	}
	return nil
}

func (uc *storage) Liveness() error {
	if err := uc.repos.Liveness(); err != nil {
		return err
	}
	return nil
}

func (uc *storage) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	uc.logger.Info("connected")
	return nil
}

func (uc *storage) Disconnect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.logger.Info("disconnected")

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *storage) Warehouse() Warehouse {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.warehose == nil {
		uc.warehose = &warehose{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.warehose
}
