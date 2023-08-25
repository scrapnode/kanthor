package storage

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

type Message interface {
	Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error)
}

type MessagePutReq struct {
	Docs []entities.Message `json:"docs" validate:"required"`
}

type MessagePutRes struct {
	Entities []entities.TSEntity
}

type message struct {
	conf    *config.Config
	logger  logging.Logger
	timer   timer.Timer
	metrics metrics.Metrics
	repos   repos.Repositories
}
