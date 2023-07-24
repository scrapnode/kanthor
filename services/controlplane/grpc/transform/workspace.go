package transform

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

func Workspace(ws *entities.Workspace) *protos.WorkspaceEntity {
	return &protos.WorkspaceEntity{
		Id:        ws.Id,
		CreatedAt: ws.CreatedAt,
		UpdatedAt: ws.UpdatedAt,
		OwnerId:   ws.OwnerId,
		Name:      ws.Name,
		Tier: &protos.WorkspaceTierEntity{
			WorkspaceId: ws.Tier.WorkspaceId,
			Name:        ws.Tier.Name,
		},
	}
}

func WorkspaceGetReq(ctx context.Context, req *protos.WorkspaceGetReq) *usecase.WorkspaceGetReq {
	return &usecase.WorkspaceGetReq{Id: req.Id}
}

func WorkspaceGetRes(ctx context.Context, res *usecase.WorkspaceGetRes) *protos.WorkspaceEntity {
	return Workspace(res.Workspace)
}
