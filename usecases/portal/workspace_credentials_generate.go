package portal

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsGenerateReq struct {
	WorkspaceId string
	Name        string
	ExpiredAt   int64
}

func (req *WorkspaceCredentialsGenerateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("workspace_id", req.WorkspaceId, entities.IdNsWs),
		validator.StringRequired("name", req.Name),
		validator.NumberGreaterThanOrEqual[int64]("expired_at", req.ExpiredAt, 0),
	)
}

type WorkspaceCredentialsGenerateRes struct {
	Credentials *entities.WorkspaceCredentials
	Password    string
}

func (uc *workspaceCredentials) Generate(ctx context.Context, req *WorkspaceCredentialsGenerateReq) (*WorkspaceCredentialsGenerateRes, error) {
	now := uc.timer.Now()
	doc := &entities.WorkspaceCredentials{
		WorkspaceId: req.WorkspaceId,
		Name:        req.Name,
	}
	doc.GenId()
	doc.SetAT(now)

	password := fmt.Sprintf("wsck_%s", utils.RandomString(constants.GlobalPasswordLength))
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
