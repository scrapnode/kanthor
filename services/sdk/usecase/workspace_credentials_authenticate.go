package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsAuthenticateReq struct {
	User string
	Pass string
}

func (req *WorkspaceCredentialsAuthenticateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("user", req.User, entities.IdNsWsc),
		validator.StringRequired("pass", req.Pass),
	)
}

type WorkspaceCredentialsAuthenticateRes struct {
	Workspace            *entities.Workspace
	WorkspaceCredentials *entities.WorkspaceCredentials
}

func (uc *workspaceCredentials) Authenticate(ctx context.Context, req *WorkspaceCredentialsAuthenticateReq) (*WorkspaceCredentialsAuthenticateRes, error) {
	key := CacheKeyWsAuthenticate(req.User)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*1, func() (*WorkspaceCredentialsAuthenticateRes, error) {
		res := &WorkspaceCredentialsAuthenticateRes{}

		credentials, err := uc.repos.WorkspaceCredentials().Get(ctx, req.User)
		if err != nil {
			return nil, err
		}
		res.WorkspaceCredentials = credentials

		expired := credentials.ExpiredAt > 0 && credentials.ExpiredAt < uc.infra.Timer.Now().UnixMilli()
		if expired {
			expiredAt := time.UnixMilli(credentials.ExpiredAt).Format(time.RFC3339)
			return nil, fmt.Errorf("workspace credentials was expired (%s)", expiredAt)
		}

		if err := uc.infra.Cryptography.KDF().StringCompare(credentials.Hash, req.Pass); err != nil {
			return nil, err
		}

		ws, err := uc.repos.Workspace().Get(ctx, credentials.WsId)
		if err != nil {
			return res, nil
		}
		res.Workspace = ws

		return res, nil
	})
}
