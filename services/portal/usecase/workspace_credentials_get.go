package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsGetReq struct {
	WsId string
	Id   string
}

func (req *WorkspaceCredentialsGetReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", req.Id, entities.IdNsWsc),
	)
}

type WorkspaceCredentialsGetRes struct {
	Doc *entities.WorkspaceCredentials
}

func (uc *workspaceCredentials) Get(ctx context.Context, req *WorkspaceCredentialsGetReq) (*WorkspaceCredentialsGetRes, error) {
	// we don't need to use cache here because the usage is too low
	wsc, err := uc.repositories.WorkspaceCredentials().Get(ctx, req.WsId, req.Id)
	if err != nil {
		return nil, err
	}

	// IMPORTANT: don't return hash value
	wsc.Hash = ""

	res := &WorkspaceCredentialsGetRes{Doc: wsc}
	return res, nil
}