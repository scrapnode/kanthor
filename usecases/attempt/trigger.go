package attempt

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/dlocker"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/attempt/repos"
)

type Trigger interface {
	Plan(ctx context.Context, req *TriggerPlanReq) (*TriggerPlanRes, error)
	Exec(ctx context.Context, req *TriggerExecReq) (*TriggerExecRes, error)
}

type trigger struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	cache     cache.Cache
	locker    dlocker.Factory
	publisher streaming.Publisher
	metrics   metric.Metrics
	repos     repos.Repositories
}
