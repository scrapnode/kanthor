package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

type EndpointRule interface {
	Create(ctx context.Context, req *EndpointRuleCreateReq) (*EndpointRuleCreateRes, error)
	Update(ctx context.Context, req *EndpointRuleUpdateReq) (*EndpointRuleUpdateRes, error)
	Delete(ctx context.Context, req *EndpointRuleDeleteReq) (*EndpointRuleDeleteRes, error)

	List(ctx context.Context, req *EndpointRuleListReq) (*EndpointRuleListRes, error)
	Get(ctx context.Context, req *EndpointRuleGetReq) (*EndpointRuleGetRes, error)
}

type EndpointRuleCreateReq struct {
	EpId string `validate:"required,startswith=ep_"`
	Name string `validate:"required"`

	Priority            int32
	Exclusionary        bool
	ConditionSource     string `validate:"required"`
	ConditionExpression string `validate:"required"`
}

type EndpointRuleCreateRes struct {
	Doc *entities.EndpointRule
}

type EndpointRuleUpdateReq struct {
	AppId string `validate:"required,startswith=app_"`
	EpId  string `validate:"required,startswith=ep_"`
	Id    string `validate:"required,startswith=epr_"`
	Name  string `validate:"required"`
}

type EndpointRuleUpdateRes struct {
	Doc *entities.EndpointRule
}

type EndpointRuleDeleteReq struct {
	AppId string `validate:"required,startswith=app_"`
	EpId  string `validate:"required,startswith=ep_"`
	Id    string `validate:"required,startswith=epr_"`
}

type EndpointRuleDeleteRes struct {
	Doc *entities.EndpointRule
}

type EndpointRuleListReq struct {
	AppId string `validate:"required,startswith=app_"`
	EpId  string `validate:"required,startswith=ep_"`
	*structure.ListReq
}

type EndpointRuleListRes struct {
	*structure.ListRes[entities.EndpointRule]
}

type EndpointRuleGetReq struct {
	AppId string `validate:"required,startswith=app_"`
	EpId  string `validate:"required,startswith=ep_"`
	Id    string `validate:"required,startswith=epr_"`
}

type EndpointRuleGetRes struct {
	Doc *entities.EndpointRule
}

type endpointRule struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	repos        repos.Repositories
}
