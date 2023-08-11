package portal

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *workspaceCredentials) Generate(ctx context.Context, req *WorkspaceCredentialsGenerateReq) (*WorkspaceCredentialsGenerateRes, error) {
	now := uc.timer.Now()
	passwords := map[string]string{}
	var docs []entities.WorkspaceCredentials
	for i := 0; i < req.Count; i++ {
		name := req.Name
		if name == "" {
			name = fmt.Sprintf("auto#%d at %s", i+1, now.Format(time.RFC3339))
		}
		credentials := entities.WorkspaceCredentials{
			WorkspaceId: req.WorkspaceId,
			Name:        name,
		}
		credentials.GenId()
		credentials.SetAT(now)

		password := utils.RandomString(constants.GlobalPasswordLength)
		passwords[credentials.Id] = password
		// once we got error, reject entirely request instead of do a partial success request
		hash, err := uc.cryptography.KDF().StringHash(password)
		if err != nil {
			return nil, err
		}

		credentials.Hash = hash
		docs = append(docs, credentials)
	}

	_, err := uc.repos.WorkspaceCredentials().BulkCreate(ctx, docs)
	if err != nil {
		return nil, err
	}

	res := &WorkspaceCredentialsGenerateRes{Credentials: docs, Passwords: passwords}
	return res, nil
}
