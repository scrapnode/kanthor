package storage

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

type Request interface {
	Put(ctx context.Context, req *RequestPutReq) (*RequestPutRes, error)
}

type RequestPutReq struct {
	Docs []entities.Request
}

type RequestPutRes struct {
	Entities []entities.Entity
}

type request struct {
	conf    *config.Config
	logger  logging.Logger
	metrics metric.Metrics
	repos   repos.Repositories
}
