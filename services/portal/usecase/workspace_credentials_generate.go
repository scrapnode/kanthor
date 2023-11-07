package usecase

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsGenerateIn struct {
	WsId      string
	Name      string
	ExpiredAt int64

	Role        string
	Permissions []authorizator.Permission
}

func (in *WorkspaceCredentialsGenerateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringRequired("name", in.Name),
		validator.NumberGreaterThanOrEqual("expired_at", in.ExpiredAt, 0),
		validator.StringRequired("role", in.Role),
		validator.SliceRequired("permissions", in.Permissions),
	)
}

type WorkspaceCredentialsGenerateOut struct {
	Credentials *entities.WorkspaceCredentials
	Password    string
}

func (uc *workspaceCredentials) Generate(ctx context.Context, in *WorkspaceCredentialsGenerateIn) (*WorkspaceCredentialsGenerateOut, error) {
	now := uc.infra.Timer.Now()
	doc := &entities.WorkspaceCredentials{
		WsId: in.WsId,
		Name: in.Name,
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

	credentials, err := uc.repositories.WorkspaceCredentials().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	if err := uc.infra.Authorizator.Grant(credentials.WsId, credentials.Id, in.Role, in.Permissions); err != nil {
		return nil, err
	}
	if err := uc.infra.Authorizator.Refresh(ctx); err != nil {
		return nil, err
	}

	res := &WorkspaceCredentialsGenerateOut{Credentials: credentials, Password: password}
	return res, nil
}
