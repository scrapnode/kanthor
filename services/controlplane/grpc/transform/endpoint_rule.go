package transform

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

func EndpointRule(epr *entities.EndpointRule) *protos.EndpointRuleEntity {
	return &protos.EndpointRuleEntity{
		Id:                  epr.Id,
		CreatedAt:           epr.CreatedAt,
		UpdatedAt:           epr.UpdatedAt,
		EndpointId:          epr.EndpointId,
		Priority:            epr.Priority,
		Exclusionary:        epr.Exclusionary,
		ConditionSource:     epr.ConditionSource,
		ConditionExpression: epr.ConditionExpression,
	}
}

func EndpointRuleListReq(ctx context.Context, req *protos.EndpointRuleListReq) *usecase.EndpointRuleListReq {
	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)
	return &usecase.EndpointRuleListReq{
		Workspace: ws,
		AppId:     req.AppId,
		EpId:      req.EpId,
		ListReq: structure.ListReq{
			Cursor: req.Cursor,
			Search: req.Search,
			Limit:  int(req.Limit),
			Ids:    req.Ids,
		},
	}
}

func EndpointRuleListRes(ctx context.Context, res *usecase.EndpointRuleListRes) *protos.EndpointRuleListRes {
	returning := &protos.EndpointRuleListRes{Cursor: res.Cursor, Data: []*protos.EndpointRuleEntity{}}
	for _, epr := range res.Data {
		returning.Data = append(returning.Data, EndpointRule(&epr))
	}
	return returning
}

func EndpointRuleGetReq(ctx context.Context, req *protos.EndpointRuleGetReq) *usecase.EndpointRuleGetReq {
	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)
	return &usecase.EndpointRuleGetReq{
		Workspace: ws,
		AppId:     req.AppId,
		EpId:      req.EpId,
		Id:        req.Id,
	}
}

func EndpointRuleGetRes(ctx context.Context, res *usecase.EndpointRuleGetRes) *protos.EndpointRuleEntity {
	return EndpointRule(res.EndpointRule)
}
