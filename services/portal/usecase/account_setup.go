package usecase

import (
	"context"
	"errors"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/project"
	"github.com/scrapnode/kanthor/services/permissions"
)

type AccountSetupIn struct {
	AccountId     string
	WorkspaceName string
}

func (in *AccountSetupIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("account_id", in.AccountId),
	)
}

type AccountSetupOut struct {
	Workspace *entities.Workspace
}

func (uc *account) Setup(ctx context.Context, in *AccountSetupIn) (*AccountSetupOut, error) {
	res, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ws, err := uc.repositories.Workspace().GetOwned(txctx, in.AccountId)
		if err != nil {
			is404 := errors.Is(err, database.ErrRecordNotFound)
			if !is404 {
				return nil, err
			}

			ws = &entities.Workspace{
				OwnerId: in.AccountId,
				Name:    in.WorkspaceName,
				Tier:    project.Tier(),
			}
			if ws.Name == "" {
				ws.Name = project.DefaultWorkspaceName()
			}

			ws.Id = suid.New(entities.IdNsWs)
			ws.SetAT(uc.infra.Timer.Now())
			if _, err := uc.repositories.Workspace().Create(ctx, ws); err != nil {
				return nil, err
			}
		}

		// check permissions
		access, err := uc.infra.Authorizator.UserPermissionsInTenant(ws.Id, in.AccountId)
		if err != nil {
			return nil, err
		}
		if len(access) == 0 {
			if err := uc.infra.Authorizator.Grant(ws.Id, in.AccountId, permissions.PortalOwner, permissions.PortalOwnerPermissions); err != nil {
				return nil, err
			}

			if err := uc.infra.Authorizator.Grant(ws.Id, in.AccountId, permissions.SdkOwner, permissions.SdkOwnerPermissions); err != nil {
				return nil, err
			}
		}

		return &AccountSetupOut{Workspace: ws}, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*AccountSetupOut), nil
}
