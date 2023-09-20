package storage

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

type Message interface {
	Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error)
}

type MessagePutReq struct {
	Docs []entities.Message
}

type MessagePutRes struct {
	Entities []entities.TSEntity
}

type message struct {
	conf    *config.Config
	logger  logging.Logger
	metrics metric.Metrics
	repos   repos.Repositories
}
