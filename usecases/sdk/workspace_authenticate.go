package sdk

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *workspace) Authenticate(ctx context.Context, req *WorkspaceAuthenticateReq) (*WorkspaceAuthenticateRes, error) {
	key := utils.Key("sdk", req.User, req.Hash)
	return cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*WorkspaceAuthenticateRes, error) {
		res := &WorkspaceAuthenticateRes{}

		credentials, err := uc.repos.WorkspaceCredentials().Get(ctx, req.User)
		if err != nil {
			return nil, err
		}
		res.WorkspaceCredentials = credentials

		expired := credentials.ExpiredAt > 0 && credentials.ExpiredAt < uc.timer.Now().UnixMilli()
		if expired {
			expiredAt := time.UnixMilli(credentials.ExpiredAt).Format(time.RFC3339)
			return nil, fmt.Errorf("workspace credentials was expired (%s)", expiredAt)
		}

		ws, err := uc.repos.Workspace().Get(ctx, credentials.WorkspaceId)
		if err != nil {
			return res, nil
		}
		res.Workspace = ws

		tier, err := uc.repos.WorkspaceTier().Get(ctx, credentials.WorkspaceId)
		if err != nil {
			return res, nil
		}
		res.WorkspaceTier = tier

		return res, nil
	})
}
