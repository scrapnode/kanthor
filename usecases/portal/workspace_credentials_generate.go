package portal

import (
	"context"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
)

func (uc *workspaceCredentials) Generate(ctx context.Context, req *WorkspaceCredentialsGenerateReq) (*WorkspaceCredentialsGenerateRes, error) {
	now := uc.timer.Now()
	doc := &entities.WorkspaceCredentials{
		WorkspaceId: req.WorkspaceId,
		Name:        req.Name,
	}
	doc.GenId()
	doc.SetAT(now)

	password := utils.RandomString(constants.GlobalPasswordLength)
	// once we got error, reject entirely request instead of do a partial success request
	hash, err := uc.cryptography.KDF().StringHash(password)
	if err != nil {
		return nil, err
	}

	doc.Hash = hash

	credentials, err := uc.repos.WorkspaceCredentials().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	res := &WorkspaceCredentialsGenerateRes{Credentials: credentials, Password: password}
	return res, nil
}
