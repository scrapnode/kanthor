package usecases

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/timer"
)

func NewDataplane(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	repos repositories.Repositories,
) Dataplane {
	logger = logger.With("usecase", "dataplane")
	return &dataplane{conf: conf, logger: logger, timer: timer, publisher: publisher, repos: repos}
}

type Dataplane interface {
	patterns.Connectable
	PutMessage(ctx context.Context, req *DataplanePutMessageReq) (*DataplanePutMessageRes, error)
}

type DataplanePutMessageReq struct {
	AppId string `json:"app_id"`
	Type  string `json:"type"`
	Body  string `json:"body"`
}

type DataplanePutMessageRes struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Bucket    string `json:"bucket"`
}

type dataplane struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	repos     repositories.Repositories
}

func (usecase *dataplane) Connect(ctx context.Context) error {
	if err := usecase.repos.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.publisher.Connect(ctx); err != nil {
		return err
	}

	usecase.logger.Info("connected")
	return nil
}

func (usecase *dataplane) Disconnect(ctx context.Context) error {
	usecase.logger.Info("disconnected")

	if err := usecase.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.publisher.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
