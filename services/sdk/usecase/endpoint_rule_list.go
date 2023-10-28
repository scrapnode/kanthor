package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleListReq struct {
	EpId string
	*structure.ListReq
}

func (req *EndpointRuleListReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ep_id", req.EpId, entities.IdNsEp),
		validator.PointerNotNil("list", req.ListReq),
	)
}

type EndpointRuleListRes struct {
	*structure.ListRes[entities.EndpointRule]
}

func (uc *endpointRule) List(ctx context.Context, req *EndpointRuleListReq) (*EndpointRuleListRes, error) {
	ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)

	ep, err := uc.repositories.Endpoint().GetOfWorkspace(ctx, ws, req.EpId)
	if err != nil {
		return nil, err
	}
	listing, err := uc.repositories.EndpointRule().List(
		ctx,
		ep,
		structure.WithListCursor(req.Cursor),
		structure.WithListSearch(req.Search),
		structure.WithListLimit(req.Limit),
		structure.WithListIds(req.Ids),
	)
	if err != nil {
		return nil, err
	}

	res := &EndpointRuleListRes{ListRes: listing}
	return res, nil
}
