package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceImportIn struct {
	Id       string
	Snapshot *entities.WorkspaceSnapshot
}

func (in *WorkspaceImportIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
		validator.PointerNotNil("snapshot", in.Snapshot),
		validator.Map(in.Snapshot.Applications, func(appId string, app entities.WorkspaceSnapshotApp) error {
			appPrefix := fmt.Sprintf("snapshot[%s]", appId)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringRequired(appPrefix+".name", app.Name),
				validator.Map(app.Endpoints, func(epId string, ep entities.WorkspaceSnapshotEp) error {
					epPrefix := fmt.Sprintf(appPrefix+".endpoints[%s]", epId)
					return validator.Validate(
						validator.DefaultConfig,
						validator.StringRequired(epPrefix+".name", ep.Name),
						validator.StringOneOf(epPrefix+".method", ep.Method, []string{http.MethodPost, http.MethodPut}),
						validator.StringUri(epPrefix+".uri", ep.Uri),
						validator.Map(ep.Rules, func(ruleId string, rule entities.WorkspaceSnapshotEpr) error {
							eprPrefix := fmt.Sprintf(epPrefix+".rules[%s]", ruleId)
							return validator.Validate(
								validator.DefaultConfig,
								validator.StringRequired(eprPrefix+".name", rule.Name),
								validator.NumberGreaterThan(eprPrefix+".priority", rule.Priority, 0),
								validator.StringRequired(eprPrefix+".condition_source", rule.ConditionSource),
								validator.StringRequired(eprPrefix+".condition_expression", rule.ConditionExpression),
							)
						}),
					)
				}),
			)
		}),
	)
}

type WorkspaceImportOut struct {
	AppIds      []string
	EpIds       []string
	EprIds      []string
	Credentials *entities.WorkspaceCredentials
}

func (uc *workspace) Import(ctx context.Context, in *WorkspaceImportIn) (*WorkspaceImportOut, error) {
	out, err := uc.repositories.Database().Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		applications := []entities.Application{}
		endpoints := []entities.Endpoint{}
		rules := []entities.EndpointRule{}

		for appId, app := range in.Snapshot.Applications {
			application := entities.Application{
				WsId: in.Id,
				Name: app.Name,
			}
			application.Id = appId
			application.SetAT(uc.infra.Timer.Now())
			applications = append(applications, application)

			for epId, ep := range app.Endpoints {
				endpoint := entities.Endpoint{
					SecretKey: utils.RandomString(32),
					AppId:     application.Id,
					Name:      ep.Name,
					Method:    ep.Method,
					Uri:       ep.Uri,
				}
				endpoint.Id = epId
				endpoint.SetAT(uc.infra.Timer.Now())
				endpoints = append(endpoints, endpoint)

				for eprId, epr := range ep.Rules {
					rule := entities.EndpointRule{
						EpId:                endpoint.Id,
						Name:                epr.Name,
						Priority:            epr.Priority,
						Exclusionary:        epr.Exclusionary,
						ConditionSource:     epr.ConditionSource,
						ConditionExpression: epr.ConditionExpression,
					}
					rule.Id = eprId
					rule.SetAT(uc.infra.Timer.Now())
					rules = append(rules, rule)
				}
			}
		}

		o := &WorkspaceImportOut{}

		appIds, err := uc.repositories.Database().Application().CreateBulk(txctx, applications)
		if err != nil {
			return nil, err
		}
		o.AppIds = appIds

		epIds, err := uc.repositories.Database().Endpoint().CreateBulk(txctx, endpoints)
		if err != nil {
			return nil, err
		}
		o.EpIds = epIds

		eprIds, err := uc.repositories.Database().EndpointRule().CreateBulk(txctx, rules)
		if err != nil {
			return nil, err
		}
		o.EprIds = eprIds

		return o, nil

	})
	if err != nil {
		return nil, err
	}

	return out.(*WorkspaceImportOut), nil
}
