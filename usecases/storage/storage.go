package storage

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

type Storage interface {
	patterns.Connectable
	Message() Message
	Request() Request
	Response() Response
}

func New(
	conf *config.Config,
	logger logging.Logger,
	metrics metric.Metrics,
	repos repos.Repositories,
) Storage {
	return &storage{
		conf:    conf,
		logger:  logger,
		metrics: metrics,
		repos:   repos,
	}
}

type storage struct {
	conf    *config.Config
	logger  logging.Logger
	metrics metric.Metrics
	repos   repos.Repositories

	mu       sync.RWMutex
	message  *message
	request  *request
	response *response
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

func (uc *storage) Message() Message {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.message == nil {
		uc.message = &message{
			conf:    uc.conf,
			logger:  uc.logger,
			metrics: uc.metrics,
			repos:   uc.repos,
		}
	}
	return uc.message
}

func (uc *storage) Request() Request {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.request == nil {
		uc.request = &request{
			conf:    uc.conf,
			logger:  uc.logger,
			metrics: uc.metrics,
			repos:   uc.repos,
		}
	}
	return uc.request
}

func (uc *storage) Response() Response {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.response == nil {
		uc.response = &response{
			conf:    uc.conf,
			logger:  uc.logger,
			metrics: uc.metrics,
			repos:   uc.repos,
		}
	}
	return uc.response
}
