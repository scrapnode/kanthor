package usecases

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/timer"
)

func NewScheduler(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	repos repositories.Repositories,
) Scheduler {
	logger = logger.With("usecase", "dataplane")
	return &scheduler{conf: conf, logger: logger, timer: timer, publisher: publisher, repos: repos}
}

type Scheduler interface {
	patterns.Connectable
	ArrangeRequests(ctx context.Context, req *ArrangeRequestsReq) (*ArrangeRequestsRes, error)
}

type ArrangeRequestsReq struct {
	Message *entities.Message
}

type ArrangeRequestsRes struct {
	Entities    []structure.BulkRes[entities.Request]
	FailKeys    []string
	SuccessKeys []string
}

type scheduler struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	repos     repositories.Repositories
}

func (service *scheduler) Connect(ctx context.Context) error {
	if err := service.repos.Connect(ctx); err != nil {
		return err
	}

	if err := service.publisher.Connect(ctx); err != nil {
		return err
	}

	service.logger.Info("connected")
	return nil
}

func (service *scheduler) Disconnect(ctx context.Context) error {
	service.logger.Info("disconnected")

	if err := service.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := service.publisher.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
