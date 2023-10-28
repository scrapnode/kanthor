package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsExpireReq struct {
	WsId     string
	Id       string
	Duration int64
}

func (req *WorkspaceCredentialsExpireReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", req.Id, entities.IdNsWsc),
		validator.NumberGreaterThanOrEqual("duration", req.Duration, 0),
	)
}

type WorkspaceCredentialsExpireRes struct {
	Id        string
	ExpiredAt int64
}

func (uc *workspaceCredentials) Expire(ctx context.Context, req *WorkspaceCredentialsExpireReq) (*WorkspaceCredentialsExpireRes, error) {
	wsc, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		wsc, err := uc.repos.WorkspaceCredentials().Get(txctx, req.WsId, req.Id)
		if err != nil {
			return nil, err
		}

		expired := wsc.ExpiredAt > 0 && wsc.ExpiredAt < uc.infra.Timer.Now().UnixMilli()
		if expired {
			return nil, errors.New("credentials was already expired")
		}

		wsc.ExpiredAt = uc.infra.Timer.Now().Add(time.Millisecond * time.Duration(req.Duration)).UnixMilli()
		wsc.SetAT(uc.infra.Timer.Now())
		return uc.repos.WorkspaceCredentials().Update(txctx, wsc)
	})
	if err != nil {
		return nil, err
	}

	doc := wsc.(*entities.WorkspaceCredentials)
	res := &WorkspaceCredentialsExpireRes{
		Id:        doc.Id,
		ExpiredAt: doc.ExpiredAt,
	}
	return res, nil
}
