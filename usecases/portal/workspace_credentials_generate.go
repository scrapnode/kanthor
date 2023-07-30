package portal

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *workspaceCredentials) Generate(ctx context.Context, req *WorkspaceCredentialsReq) (*WorkspaceCredentialsRes, error) {
	now := uc.timer.Now()
	var docs []entities.WorkspaceCredentials
	for i := 0; i < req.Count; i++ {
		credentials := entities.WorkspaceCredentials{WorkspaceId: req.WorkspaceId}
		credentials.GenId()
		credentials.SetAT(now)
		credentials.SetDefaultExpired(now)
		docs = append(docs, credentials)
	}

	_, err := uc.repos.WorkspaceCredentials().BulkCreate(ctx, docs)
	if err != nil {
		return nil, err
	}

	res := &WorkspaceCredentialsRes{Credentials: docs}
	return res, nil
}
