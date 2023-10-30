package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
)

type Endeavor interface {
	Plan(ctx context.Context, req *EndeavorPlanReq) (*EndeavorPlanRes, error)
	Exec(ctx context.Context, req *EndeavorExecReq) (*EndeavorExecRes, error)
}

type endeavor struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
	publisher    streaming.Publisher
}
