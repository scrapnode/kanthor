package portal

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsGetReq struct {
	WorkspaceId string
	Id          string
}

func (req *WorkspaceCredentialsGetReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("workspace_id", req.WorkspaceId, "ws_"),
		validator.StringStartsWith("id", req.Id, "wsc_"),
	)
}

type WorkspaceCredentialsGetRes struct {
	Doc *entities.WorkspaceCredentials
}

func (uc *workspaceCredentials) Get(ctx context.Context, req *WorkspaceCredentialsGetReq) (*WorkspaceCredentialsGetRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	// we don't need to use cache here because the usage is too low
	wsc, err := uc.repos.WorkspaceCredentials().Get(ctx, ws.Id, req.Id)
	if err != nil {
		return nil, err
	}

	// IMPORTANT: don't return hash value
	wsc.Hash = ""

	res := &WorkspaceCredentialsGetRes{Doc: wsc}
	return res, nil
}
