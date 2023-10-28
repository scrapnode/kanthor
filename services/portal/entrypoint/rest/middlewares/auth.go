package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

var (
	HeaderWorkspace = "kanthor-ws-id"
)

func UseAuth(engine authenticator.Authenticator, uc usecase.Portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)

		// authenticate
		credentials := ginctx.Request.Header.Get(authenticator.HeaderAuth)
		acc, err := engine.Verify(credentials)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx = context.WithValue(ctx, authenticator.CtxAcc, acc)

		ws, err := workspace(uc, ctx, ginctx.Request.Header.Get(HeaderWorkspace))
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx = context.WithValue(ctx, gateway.CtxWs, ws)

		// authorize ownership
		isOwner := acc.Sub == ws.OwnerId
		ctx = context.WithValue(ctx, gateway.CtxAuhzOk, isOwner)

		ginctx.Set(gateway.KeyContext, ctx)
		ginctx.Next()
	}
}

func workspace(uc usecase.Portal, ctx context.Context, id string) (*entities.Workspace, error) {
	req := &usecase.WorkspaceGetReq{Id: id}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	res, err := uc.Workspace().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.Workspace, nil
}
