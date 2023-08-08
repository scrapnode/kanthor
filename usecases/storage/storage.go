package storage

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
	"sync"
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
	repos repos.Repositories,
) Storage {
	return &storage{
		conf:   conf,
		logger: logger,
		repos:  repos,
	}
}

type storage struct {
	conf   *config.Config
	logger logging.Logger
	repos  repos.Repositories

	mu       sync.RWMutex
	message  *message
	request  *request
	response *response
}

func (uc *storage) Connect(ctx context.Context) error {
	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	uc.logger.Info("connected")
	return nil
}

func (uc *storage) Disconnect(ctx context.Context) error {
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
			conf:   uc.conf,
			logger: uc.logger,
			repos:  uc.repos,
		}
	}
	return uc.message
}

func (uc *storage) Request() Request {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.request == nil {
		uc.request = &request{
			conf:   uc.conf,
			logger: uc.logger,
			repos:  uc.repos,
		}
	}
	return uc.request
}

func (uc *storage) Response() Response {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.response == nil {
		uc.response = &response{
			conf:   uc.conf,
			logger: uc.logger,
			repos:  uc.repos,
		}
	}
	return uc.response
}
