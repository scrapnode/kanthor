package sdk

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceCredentialsExpireReq struct {
	User      string
	ExpiredAt int64
}

func (req *WorkspaceCredentialsExpireReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("user", req.User, entities.IdNsWsc),
		validator.NumberGreaterThan("expired_at", req.ExpiredAt, 0),
	)
}

type WorkspaceCredentialsExpireRes struct {
	Ok bool
}

func (uc *workspaceCredentials) Expire(ctx context.Context, req *WorkspaceCredentialsExpireReq) (*WorkspaceCredentialsExpireRes, error) {
	key := CacheKeyWsAuthenticate(req.User)
	at := time.UnixMilli(req.ExpiredAt)
	ok, err := uc.cache.ExpireAt(ctx, key, at)
	if err != nil {
		return nil, err
	}

	res := &WorkspaceCredentialsExpireRes{Ok: ok}
	return res, nil
}
