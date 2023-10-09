package attempt

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/attempt/repos"
)

type Application interface {
	Trigger(ctx context.Context, req *ApplicationTriggerReq) (*ApplicationTriggerRes, error)
}

type application struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	cache     cache.Cache
	publisher streaming.Publisher
	metrics   metric.Metrics
	repos     repos.Repositories
}
