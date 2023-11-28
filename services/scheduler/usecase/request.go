package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/services/scheduler/config"
	"github.com/scrapnode/kanthor/services/scheduler/repositories"
)

type Request interface {
	Schedule(ctx context.Context, in *RequestScheduleIn) (*RequestScheduleOut, error)
}

type request struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	publisher    streaming.Publisher
	repositories repositories.Repositories
	duplicated   *safe.Map[int]
}
