package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/internal/assessor"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
)

type Trigger interface {
	Perform(
		ctx context.Context,
		appId string,
		msgs map[string]entities.Message,
		applicable *assessor.Assets,
		attemptDelay int64,
	) (*TriggerExecOut, error)
	Applicable(ctx context.Context, appId string) (*assessor.Assets, error)

	Plan(ctx context.Context, in *TriggerPlanIn) (*TriggerPlanOut, error)
	Exec(ctx context.Context, in *TriggerExecIn) (*TriggerExecOut, error)
}

type trigger struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
