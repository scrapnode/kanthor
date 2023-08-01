package sdk

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *workspace) Authenticate(ctx context.Context, req *WorkspaceAuthenticateReq) (*WorkspaceAuthenticateRes, error) {
	key := utils.Key("sdk", "workspace", req.User, req.Hash)
	// No mather the function got error or not, we want to cache the result of authentication for next usage,
	// so we MUST NOT return error
	return cache.Warp[WorkspaceAuthenticateRes](uc.cache, key, time.Hour*24, func() (*WorkspaceAuthenticateRes, error) {
		res := &WorkspaceAuthenticateRes{}

		credentials, err := uc.repos.WorkspaceCredentials().Get(ctx, req.User)
		if err != nil {
			res.Error = err
			return res, nil
		}

		expired := credentials.ExpiredAt > 0 && credentials.ExpiredAt < uc.timer.Now().UnixMilli()
		if expired {
			expiredAt := time.UnixMilli(credentials.ExpiredAt).Format(time.RFC3339)
			res.Error = fmt.Errorf("workspace credentials was expired (%s)", expiredAt)
			return res, nil
		}

		ws, err := uc.repos.Workspace().Get(ctx, credentials.WorkspaceId)
		if err != nil {
			res.Error = err
			return res, nil
		}
		res.Workspace = ws

		tier, err := uc.repos.WorkspaceTier().Get(ctx, credentials.WorkspaceId)
		if err != nil {
			res.Error = err
			return res, nil
		}
		res.WorkspaceTier = tier

		return res, nil
	})
}
