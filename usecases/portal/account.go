package portal

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/usecases/portal/repos"
)

type Account interface {
	Setup(ctx context.Context, req *AccountSetupReq) (*AccountSetupRes, error)
}

type account struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories
}
