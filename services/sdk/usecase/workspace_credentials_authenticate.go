package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsAuthenticateIn struct {
	User string
	Pass string
}

func (in *WorkspaceCredentialsAuthenticateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("user", in.User, entities.IdNsWsc),
		validator.StringRequired("pass", in.Pass),
	)
}

type WorkspaceCredentialsAuthenticateOut struct {
	Workspace            *entities.Workspace
	WorkspaceCredentials *entities.WorkspaceCredentials
}

func (uc *workspaceCredentials) Authenticate(ctx context.Context, in *WorkspaceCredentialsAuthenticateIn) (*WorkspaceCredentialsAuthenticateOut, error) {
	key := CacheKeyWsAuthenticate(in.User)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*1, func() (*WorkspaceCredentialsAuthenticateOut, error) {
		res := &WorkspaceCredentialsAuthenticateOut{}

		credentials, err := uc.repositories.WorkspaceCredentials().Get(ctx, in.User)
		if err != nil {
			return nil, err
		}
		res.WorkspaceCredentials = credentials

		expired := credentials.ExpiredAt > 0 && credentials.ExpiredAt < uc.infra.Timer.Now().UnixMilli()
		if expired {
			expiredAt := time.UnixMilli(credentials.ExpiredAt).Format(time.RFC3339)
			return nil, fmt.Errorf("workspace credentials was expired (%s)", expiredAt)
		}

		if err := uc.infra.Cryptography.KDF().StringCompare(credentials.Hash, in.Pass); err != nil {
			return nil, err
		}

		ws, err := uc.repositories.Workspace().Get(ctx, credentials.WsId)
		if err != nil {
			return res, nil
		}
		res.Workspace = ws

		return res, nil
	})
}
