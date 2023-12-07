package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/structure"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleListIn struct {
	Ws   *entities.Workspace
	EpId string
	*structure.ListReq
}

func (in *EndpointRuleListIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.PointerNotNil("ws", in.Ws),
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
		validator.PointerNotNil("list", in.ListReq),
	)
}

type EndpointRuleListOut struct {
	*structure.ListRes[entities.EndpointRule]
}

func (uc *endpointRule) List(ctx context.Context, in *EndpointRuleListIn) (*EndpointRuleListOut, error) {
	ep, err := uc.repositories.Endpoint().GetOfWorkspace(ctx, in.Ws, in.EpId)
	if err != nil {
		return nil, err
	}
	listing, err := uc.repositories.EndpointRule().List(
		ctx,
		ep,
		structure.WithListCursor(in.Cursor),
		structure.WithListSearch(in.Search),
		structure.WithListLimit(in.Limit),
		structure.WithListIds(in.Ids),
	)
	if err != nil {
		return nil, err
	}

	return &EndpointRuleListOut{ListRes: listing}, nil
}
