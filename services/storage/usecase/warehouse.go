package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/storage/config"
	"github.com/scrapnode/kanthor/services/storage/repositories"
)

type Warehouse interface {
	Put(ctx context.Context, req *WarehousePutReq) (*WarehousePutRes, error)
}

type warehose struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
