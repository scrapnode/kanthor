package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
)

type Retry interface {
	Trigger(ctx context.Context, in *RetryTriggerIn) (*RetryTriggerOut, error)
	Select(ctx context.Context, in *RetrySelectIn) (*RetrySelectOut, error)
	Endeavor(ctx context.Context, in *RetryEndeavorIn) (*RetryEndeavorOut, error)
}

type retry struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	publisher    streaming.Publisher
	repositories repositories.Repositories
}
