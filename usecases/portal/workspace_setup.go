package portal

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceSetupReq struct {
	Workspace     *entities.Workspace
	Applications  []entities.Application
	Endpoints     []entities.Endpoint
	EndpointRules []entities.EndpointRule
}

func (req *WorkspaceSetupReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.PointerNotNil("wokrspace", req.Workspace),
		validator.Array(req.Applications, func(i int, item *entities.Application) error {
			prefix := fmt.Sprintf("req.applications[%d]", i)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringStartsWith(prefix+".id", item.Id, entities.IdNsApp),
				validator.NumberGreaterThan(prefix+".created_at", item.CreatedAt, 0),
				validator.NumberGreaterThan(prefix+".updated_at", item.UpdatedAt, 0),
				validator.StringStartsWith(prefix+".workspace_id", item.WsId, entities.IdNsWs),
				validator.StringRequired(prefix+".name", item.Name),
			)
		}),
		validator.Array(req.Endpoints, func(i int, item *entities.Endpoint) error {
			prefix := fmt.Sprintf("req.endpoints[%d]", i)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringStartsWith(prefix+".id", item.Id, entities.IdNsEp),
				validator.NumberGreaterThan(prefix+".created_at", item.CreatedAt, 0),
				validator.NumberGreaterThan(prefix+".updated_at", item.UpdatedAt, 0),
				validator.StringStartsWith(prefix+".app_id", item.AppId, entities.IdNsApp),
				validator.StringRequired(prefix+".name", item.Name),
				validator.StringRequired(prefix+".secret_key", item.SecretKey),
				validator.StringRequired(prefix+".method", item.Method),
				validator.StringUri(prefix+".uri", item.Uri),
			)
		}),
		validator.Array(req.EndpointRules, func(i int, item *entities.EndpointRule) error {
			prefix := fmt.Sprintf("req.endpoint_rules[%d]", i)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringStartsWith(prefix+".id", item.Id, entities.IdNsEpr),
				validator.NumberGreaterThan(prefix+".created_at", item.CreatedAt, 0),
				validator.NumberGreaterThan(prefix+".updated_at", item.UpdatedAt, 0),
				validator.StringStartsWith(prefix+".endpoint_id", item.EpId, entities.IdNsEp),
				validator.StringRequired(prefix+".name", item.Name),
				validator.NumberGreaterThanOrEqual(prefix+".priority", item.Priority, 0),
				validator.StringRequired(prefix+".condition_source", item.ConditionSource),
				validator.StringRequired(prefix+".condition_expression", item.ConditionExpression),
			)
		}),
	)
}

type WorkspaceSetupRes struct {
	ApplicationIds  []string
	EndpointIds     []string
	EndpointRuleIds []string
	Status          map[string]bool
}

func (uc *workspace) Setup(ctx context.Context, req *WorkspaceSetupReq) (*WorkspaceSetupRes, error) {
	res, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		// starting with false
		status := map[string]bool{}

		for _, app := range req.Applications {
			status[app.Id] = false
		}
		for _, ep := range req.Endpoints {
			status[ep.Id] = false
		}
		for _, epr := range req.EndpointRules {
			status[epr.Id] = false
		}

		appIds, err := uc.repos.Application().BulkCreate(txctx, req.Applications)
		if err != nil {
			return nil, err
		}
		for _, appId := range appIds {
			status[appId] = true
		}

		epIds, err := uc.repos.Endpoint().BulkCreate(txctx, req.Endpoints)
		if err != nil {
			return nil, err
		}
		for _, epId := range epIds {
			status[epId] = true
		}

		eprIds, err := uc.repos.EndpointRule().BulkCreate(txctx, req.EndpointRules)
		if err != nil {
			return nil, err
		}
		for _, eprId := range eprIds {
			status[eprId] = true
		}

		res := &WorkspaceSetupRes{
			ApplicationIds:  appIds,
			EndpointIds:     epIds,
			EndpointRuleIds: eprIds,
			Status:          status,
		}
		return res, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*WorkspaceSetupRes), nil
}
