package storage

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

type Response interface {
	Put(ctx context.Context, req *ResponsePutReq) (*ResponsePutRes, error)
}

type ResponsePutReq struct {
	Docs []entities.Response
}

type ResponsePutRes struct {
	Entities []entities.Entity
}

type response struct {
	conf    *config.Config
	logger  logging.Logger
	metrics metric.Metrics
	repos   repos.Repositories
}
