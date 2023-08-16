package sdk

import (
	"context"
	"time"
)

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
