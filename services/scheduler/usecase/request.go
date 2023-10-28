package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services/scheduler/config"
	"github.com/scrapnode/kanthor/services/scheduler/repos"
)

type Request interface {
	Schedule(ctx context.Context, req *RequestScheduleReq) (*RequestScheduleRes, error)
}

type request struct {
	conf      *config.Config
	logger    logging.Logger
	infra     *infrastructure.Infrastructure
	publisher streaming.Publisher
	repos     repos.Repositories
}
