package attempt

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/usecases/attempt/repos"
)

type Trigger interface {
	Plan(ctx context.Context, req *TriggerPlanReq) (*TriggerPlanRes, error)
	Exec(ctx context.Context, req *TriggerExecReq) (*TriggerExecRes, error)
}

type trigger struct {
	conf      *config.Config
	logger    logging.Logger
	infra     *infrastructure.Infrastructure
	publisher streaming.Publisher
	repos     repos.Repositories
}
