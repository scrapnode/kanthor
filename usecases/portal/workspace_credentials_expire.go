package portal

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"time"
)

func (uc *workspaceCredentials) Expire(ctx context.Context, req *WorkspaceCredentialsExpireReq) (*WorkspaceCredentialsExpireRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	wsc, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		wsc, err := uc.repos.WorkspaceCredentials().Get(txctx, ws.Id, req.Id)
		if err != nil {
			return nil, err
		}

		expired := wsc.ExpiredAt > 0 && wsc.ExpiredAt < uc.timer.Now().UnixMilli()
		if expired {
			return nil, errors.New("credentials was already expired")
		}

		wsc.ExpiredAt = uc.timer.Now().Add(time.Millisecond * time.Duration(req.Duration)).UnixMilli()
		wsc.SetAT(uc.timer.Now())
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
