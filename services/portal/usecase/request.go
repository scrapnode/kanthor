package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repositories"
)

type Request interface {
	List(ctx context.Context, in *RequestListIn) (*RequestListOut, error)
	Get(ctx context.Context, in *RequestGetIn) (*RequestGetOut, error)
}

type request struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
