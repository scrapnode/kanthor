package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
)

func UseAuthz(authz authorizator.Authorizator, uc portaluc.Portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		wsId := ginctx.Request.Header.Get(authorizator.HeaderWorkspace)

		req := &portaluc.WorkspaceGetReq{Id: wsId}
		if err := req.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := uc.Workspace().Get(ctx, req)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		acc := ctx.Value(authenticator.CtxAcc).(*authenticator.Account)
		isOwner := acc.Sub == res.Workspace.OwnerId
		// if account is not owned this workspace, authorize the request
		if !isOwner {
			obj := ginctx.FullPath() // form of /application/:app_id
			act := ginctx.Request.Method
			ok, err := authz.Enforce(res.Workspace.Id, acc.Sub, obj, act)
			if err != nil {
				ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if !ok {
				ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you have no permission to perform the action"})
				return
			}
		}

		ctx = context.WithValue(ctx, authorizator.CtxWs, res.Workspace)
		ginctx.Set(gateway.KeyCtx, ctx)
		ginctx.Next()
	}
}
