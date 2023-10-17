package portal

import (
	"context"
	"errors"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsExpireReq struct {
	WorkspaceId string
	Id          string
	Duration    int64
}

func (req *WorkspaceCredentialsExpireReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("workspace_id", req.WorkspaceId, entities.IdNsWs),
		validator.StringStartsWith("id", req.Id, entities.IdNsWsc),
		validator.NumberGreaterThanOrEqual[int64]("duration", req.Duration, 0),
	)
}

type WorkspaceCredentialsExpireRes struct {
	Id        string
	ExpiredAt int64
}

func (uc *workspaceCredentials) Expire(ctx context.Context, req *WorkspaceCredentialsExpireReq) (*WorkspaceCredentialsExpireRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	wsc, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		wsc, err := uc.repos.WorkspaceCredentials().Get(txctx, ws.Id, req.Id)
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
