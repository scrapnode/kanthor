package storage

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/usecases/storage/repos"
)

type Warehouse interface {
	Put(ctx context.Context, req *WarehousePutReq) (*WarehousePutRes, error)
}

type warehose struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories
}
