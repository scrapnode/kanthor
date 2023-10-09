package scheduler

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

type Request interface {
	Arrange(ctx context.Context, req *RequestArrangeReq) (*RequestArrangeRes, error)
	Schedule(ctx context.Context, req *RequestScheduleReq) (*RequestScheduleRes, error)
}

type request struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	cache     cache.Cache
	metrics   metric.Metrics
	repos     repos.Repositories
}
