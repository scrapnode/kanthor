package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsExpireIn struct {
	WsId     string
	Id       string
	Duration int64
}

func (in *WorkspaceCredentialsExpireIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsWsc),
		validator.NumberGreaterThanOrEqual("duration", in.Duration, 0),
	)
}

type WorkspaceCredentialsExpireOut struct {
	Id        string
	ExpiredAt int64
}

func (uc *workspaceCredentials) Expire(ctx context.Context, in *WorkspaceCredentialsExpireIn) (*WorkspaceCredentialsExpireOut, error) {
	wsc, err := uc.repositories.Database().Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		wsc, err := uc.repositories.Database().WorkspaceCredentials().Get(txctx, in.WsId, in.Id)
		if err != nil {
			return nil, err
		}

		expired := wsc.ExpiredAt > 0 && wsc.ExpiredAt < uc.infra.Timer.Now().UnixMilli()
		if expired {
			return nil, errors.New("credentials was already expired")
		}

		expiredAt := uc.infra.Timer.Now().Add(time.Millisecond * time.Duration(in.Duration)).UnixMilli()
		if wsc.ExpiredAt > 0 && expiredAt > wsc.ExpiredAt {
			return nil, errors.New("credentials expired could not be extended with longer expire time")
		}

		wsc.ExpiredAt = expiredAt
		wsc.SetAT(uc.infra.Timer.Now())
		return uc.repositories.Database().WorkspaceCredentials().Update(txctx, wsc)
	})
	if err != nil {
		return nil, err
	}

	doc := wsc.(*entities.WorkspaceCredentials)
	res := &WorkspaceCredentialsExpireOut{
		Id:        doc.Id,
		ExpiredAt: doc.ExpiredAt,
	}
	return res, nil
}
