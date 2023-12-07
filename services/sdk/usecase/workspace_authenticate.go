package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceAuthenticateIn struct {
	User string
	Pass string
}

func (in *WorkspaceAuthenticateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("user", in.User, entities.IdNsWsc),
		validator.StringRequired("pass", in.Pass),
	)
}

type WorkspaceAuthenticateOut struct {
	Credentials *entities.WorkspaceCredentials
}

func (uc *workspace) Authenticate(ctx context.Context, in *WorkspaceAuthenticateIn) (*WorkspaceAuthenticateOut, error) {
	key := CacheKeyWsAuthenticate(in.User)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*1, func() (*WorkspaceAuthenticateOut, error) {
		credentials, err := uc.repositories.WorkspaceCredentials().Get(ctx, in.User)
		if err != nil {
			return nil, err
		}

		expired := credentials.ExpiredAt > 0 && credentials.ExpiredAt < uc.infra.Timer.Now().UnixMilli()
		if expired {
			expiredAt := time.UnixMilli(credentials.ExpiredAt).Format(time.RFC3339)
			return nil, fmt.Errorf("workspace credentials was expired (%s)", expiredAt)
		}

		if err := utils.PasswordCompare(in.Pass, credentials.Hash); err != nil {
			return nil, err
		}

		return &WorkspaceAuthenticateOut{Credentials: credentials}, nil
	})
}
