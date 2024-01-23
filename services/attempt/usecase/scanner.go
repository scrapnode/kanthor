package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
)

type Scanner interface {
	Schedule(ctx context.Context, in *ScannerScheduleIn) (*ScannerScheduleOut, error)
	Execute(ctx context.Context, in *ScannerExecuteIn) (*ScannerExecuteOut, error)
}

type scanner struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	publisher    streaming.Publisher
	repositories repositories.Repositories
}
