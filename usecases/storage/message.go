package storage

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

type Message interface {
	Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error)
}

type message struct {
	conf    *config.Config
	logger  logging.Logger
	metrics metric.Metrics
	repos   repos.Repositories
}
