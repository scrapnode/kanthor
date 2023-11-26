package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
)

type Cli interface {
	TriggerExecWithDateRange(ctx context.Context, in *TriggerExecWithDateRangeIn) (*TriggerExecOut, error)
	TriggerExecWithMessageIds(ctx context.Context, in *TriggerExecWithMessageIdsIn) (*TriggerExecOut, error)
}

type cli struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories

	trigger  Trigger
	endeavor Endeavor
}
