package message

import (
	"context"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/infrastructure/auth"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Service interface {
	patterns.Connectable
	Create(ctx context.Context, req *CreateReq) (*CreateRes, error)
}

func NewService(
	conf *config.Config,
	logger logging.Logger,
	auth auth.Auth,
	publisher streaming.Publisher,
	repo Repository,
) Service {
	logger = logger.With("component", "dataplane.usecases.message")
	return &service{conf: conf, logger: logger, auth: auth, publisher: publisher, repo: repo}
}

type service struct {
	conf      *config.Config
	logger    logging.Logger
	auth      auth.Auth
	publisher streaming.Publisher
	repo      Repository
}

func (service *service) Connect(ctx context.Context) error {
	if err := service.repo.Connect(ctx); err != nil {
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

	if err := service.repo.Disconnect(ctx); err != nil {
		return err
	}

	if err := service.publisher.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
