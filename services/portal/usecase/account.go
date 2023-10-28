package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repositories"
)

type Account interface {
	Setup(ctx context.Context, req *AccountSetupReq) (*AccountSetupRes, error)
}

type account struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
