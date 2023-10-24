package portal

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsGenerateReq struct {
	WsId      string
	Name      string
	ExpiredAt int64

	Role        string
	Permissions []authorizator.Permission
}

func (req *WorkspaceCredentialsGenerateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringRequired("name", req.Name),
		validator.NumberGreaterThanOrEqual("expired_at", req.ExpiredAt, 0),
		validator.StringRequired("role", req.Role),
		validator.SliceRequired("permissions", req.Permissions),
	)
}

type WorkspaceCredentialsGenerateRes struct {
	Credentials *entities.WorkspaceCredentials
	Password    string
}

func (uc *workspaceCredentials) Generate(ctx context.Context, req *WorkspaceCredentialsGenerateReq) (*WorkspaceCredentialsGenerateRes, error) {
	now := uc.infra.Timer.Now()
	doc := &entities.WorkspaceCredentials{
		WsId: req.WsId,
		Name: req.Name,
	}
	doc.GenId()
	doc.SetAT(now)

	password := fmt.Sprintf("wsck_%s", utils.RandomString(constants.PasswordLength))
	// once we got error, reject entirely request instead of do a partial success request
	hash, err := uc.infra.Cryptography.KDF().StringHash(password)
	if err != nil {
		return nil, err
	}
	doc.Hash = hash

	credentials, err := uc.repos.WorkspaceCredentials().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	if err := uc.infra.Authorizator.Grant(credentials.WsId, credentials.Id, req.Role, req.Permissions); err != nil {
		return nil, err
	}
	if err := uc.infra.Authorizator.Refresh(ctx); err != nil {
		return nil, err
	}

	res := &WorkspaceCredentialsGenerateRes{Credentials: credentials, Password: password}
	return res, nil
}
