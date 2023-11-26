package usecase

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceSetupIn struct {
	Workspace     *entities.Workspace
	Applications  []entities.Application
	Endpoints     []entities.Endpoint
	EndpointRules []entities.EndpointRule
}

func (in *WorkspaceSetupIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.PointerNotNil("wokrspace", in.Workspace),
		validator.Slice(in.Applications, func(i int, item *entities.Application) error {
			prefix := fmt.Sprintf("in.applications[%d]", i)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringStartsWith(prefix+".id", item.Id, entities.IdNsApp),
				validator.NumberGreaterThan(prefix+".created_at", item.CreatedAt, 0),
				validator.NumberGreaterThan(prefix+".updated_at", item.UpdatedAt, 0),
				validator.StringStartsWith(prefix+".ws_id", item.WsId, entities.IdNsWs),
				validator.StringRequired(prefix+".name", item.Name),
			)
		}),
		validator.Slice(in.Endpoints, func(i int, item *entities.Endpoint) error {
			prefix := fmt.Sprintf("in.endpoints[%d]", i)
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
		validator.Slice(in.EndpointRules, func(i int, item *entities.EndpointRule) error {
			prefix := fmt.Sprintf("in.endpoint_rules[%d]", i)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringStartsWith(prefix+".id", item.Id, entities.IdNsEpr),
				validator.NumberGreaterThan(prefix+".created_at", item.CreatedAt, 0),
				validator.NumberGreaterThan(prefix+".updated_at", item.UpdatedAt, 0),
				validator.StringStartsWith(prefix+".ep_id", item.EpId, entities.IdNsEp),
				validator.StringRequired(prefix+".name", item.Name),
				validator.NumberGreaterThanOrEqual(prefix+".priority", item.Priority, 0),
				validator.StringRequired(prefix+".condition_source", item.ConditionSource),
				validator.StringRequired(prefix+".condition_expression", item.ConditionExpression),
			)
		}),
	)
}

type WorkspaceSetupOut struct {
	ApplicationIds  []string
	EndpointIds     []string
	EndpointRuleIds []string
	Status          map[string]bool
}

func (uc *workspace) Setup(ctx context.Context, in *WorkspaceSetupIn) (*WorkspaceSetupOut, error) {
	res, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		// starting with false
		status := map[string]bool{}

		for _, app := range in.Applications {
			status[app.Id] = false
		}
		for _, ep := range in.Endpoints {
			status[ep.Id] = false
		}
		for _, epr := range in.EndpointRules {
			status[epr.Id] = false
		}

		appIds, err := uc.repositories.Application().BulkCreate(txctx, in.Applications)
		if err != nil {
			return nil, err
		}
		for _, appId := range appIds {
			status[appId] = true
		}

		epIds, err := uc.repositories.Endpoint().BulkCreate(txctx, in.Endpoints)
		if err != nil {
			return nil, err
		}
		for _, epId := range epIds {
			status[epId] = true
		}

		eprIds, err := uc.repositories.EndpointRule().BulkCreate(txctx, in.EndpointRules)
		if err != nil {
			return nil, err
		}
		for _, eprId := range eprIds {
			status[eprId] = true
		}

		res := &WorkspaceSetupOut{
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
	return res.(*WorkspaceSetupOut), nil
}
