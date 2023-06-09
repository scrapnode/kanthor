package message

import (
	"context"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/timer"
)

type Service interface {
	patterns.Connectable
	Create(ctx context.Context, req *CreateReq) (*CreateRes, error)
}

func NewService(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	repos repositories.Repositories,
) Service {
	logger = logger.With("component", "dataplane.usecases.message")
	return &service{conf: conf, logger: logger, timer: timer, publisher: publisher, repos: repos}
}

type service struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	repos     repositories.Repositories
}

func (service *service) Connect(ctx context.Context) error {
	if err := service.repos.Connect(ctx); err != nil {
		return err
	}

	if err := service.publisher.Connect(ctx); err != nil {
		return err
	}

	service.logger.Info("connected")
	return nil
}

func (service *service) Disconnect(ctx context.Context) error {
	service.logger.Info("disconnected")

	if err := service.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := service.publisher.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
