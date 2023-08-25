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

type Response interface {
	Put(ctx context.Context, req *ResponsePutReq) (*ResponsePutRes, error)
}

type ResponsePutReq struct {
	Docs []entities.Response `json:"docs" validate:"required"`
}

type ResponsePutRes struct {
	Entities []entities.TSEntity
}

type response struct {
	conf    *config.Config
	logger  logging.Logger
	timer   timer.Timer
	metrics metrics.Metrics
	repos   repos.Repositories
}
