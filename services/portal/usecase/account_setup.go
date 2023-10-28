package usecase

import (
	"context"
	"errors"

	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/project"
)

type AccountSetupReq struct {
	AccountId string
}

func (req *AccountSetupReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("account_id", req.AccountId),
	)
}

type AccountSetupRes struct {
	Workspace *entities.Workspace
}

func (uc *account) Setup(ctx context.Context, req *AccountSetupReq) (*AccountSetupRes, error) {
	res, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ws, err := uc.repositories.Workspace().GetOwned(txctx, req.AccountId)
		if err != nil {
			is404 := errors.Is(err, database.ErrRecordNotFound)
			if !is404 {
				return nil, err
			}

			ws = &entities.Workspace{OwnerId: req.AccountId, Name: project.DefaultWorkspaceName(), Tier: project.Tier()}
			ws.GenId()
			ws.SetAT(uc.infra.Timer.Now())
			if _, err := uc.repositories.Workspace().Create(ctx, ws); err != nil {
				return nil, err
			}
		}

		return &AccountSetupRes{Workspace: ws}, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*AccountSetupRes), nil
}
