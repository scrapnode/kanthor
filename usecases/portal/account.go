package portal

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type Account interface {
	Setup(ctx context.Context, req *AccountSetupReq) (*AccountSetupRes, error)
}

type AccountSetupReq struct {
	AccountId string `json:"account_id"`
}

type AccountSetupRes struct {
	Workspace     *entities.Workspace     `json:"workspace"`
	WorkspaceTier *entities.WorkspaceTier `json:"workspace_tier"`
}

type account struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	metrics      metric.Metrics
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
