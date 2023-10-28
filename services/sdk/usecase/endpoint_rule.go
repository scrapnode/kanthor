package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/repos"
)

type EndpointRule interface {
	Create(ctx context.Context, req *EndpointRuleCreateReq) (*EndpointRuleCreateRes, error)
	Update(ctx context.Context, req *EndpointRuleUpdateReq) (*EndpointRuleUpdateRes, error)
	Delete(ctx context.Context, req *EndpointRuleDeleteReq) (*EndpointRuleDeleteRes, error)

	List(ctx context.Context, req *EndpointRuleListReq) (*EndpointRuleListRes, error)
	Get(ctx context.Context, req *EndpointRuleGetReq) (*EndpointRuleGetRes, error)
}

type endpointRule struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	repos  repos.Repositories
}
