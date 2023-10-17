package portal

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsUpdateReq struct {
	WorkspaceId string
	Id          string
	Name        string
}

func (req *WorkspaceCredentialsUpdateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("workspace_id", req.WorkspaceId, entities.IdNsWs),
		validator.StringStartsWith("id", req.Id, entities.IdNsWsc),
		validator.StringRequired("name", req.Name),
	)
}

type WorkspaceCredentialsUpdateRes struct {
	Doc *entities.WorkspaceCredentials
}

func (uc *workspaceCredentials) Update(ctx context.Context, req *WorkspaceCredentialsUpdateReq) (*WorkspaceCredentialsUpdateRes, error) {
	doc, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		wsc, err := uc.repos.WorkspaceCredentials().Get(txctx, req.WorkspaceId, req.Id)
		if err != nil {
			return nil, err
		}

		wsc.Name = req.Name
		wsc.SetAT(uc.infra.Timer.Now())
		return uc.repos.WorkspaceCredentials().Update(txctx, wsc)
	})
	if err != nil {
		return nil, err
	}

	wsc := doc.(*entities.WorkspaceCredentials)
	// IMPORTANT: don't return hash value
	wsc.Hash = ""

	res := &WorkspaceCredentialsUpdateRes{Doc: wsc}
	return res, nil
}
