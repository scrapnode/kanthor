package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/repositories"
)

type EndpointRule interface {
	Create(ctx context.Context, in *EndpointRuleCreateIn) (*EndpointRuleCreateOut, error)
	Update(ctx context.Context, in *EndpointRuleUpdateIn) (*EndpointRuleUpdateOut, error)
	Delete(ctx context.Context, in *EndpointRuleDeleteIn) (*EndpointRuleDeleteOut, error)

	List(ctx context.Context, in *EndpointRuleListIn) (*EndpointRuleListOut, error)
	Get(ctx context.Context, in *EndpointRuleGetIn) (*EndpointRuleGetOut, error)
}

type endpointRule struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
