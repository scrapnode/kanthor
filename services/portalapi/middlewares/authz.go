package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
	"net/http"
)

func UseAuthz(authz authorizator.Authorizator, uc portaluc.Portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet("ctx").(context.Context)
		wsId := ginctx.Request.Header.Get(authorizator.HeaderWorkspace)
		res, err := uc.Workspace().Get(ctx, &portaluc.WorkspaceGetReq{Id: wsId})
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		acc := ctx.Value(authenticator.CtxAcc).(*authenticator.Account)

		obj := ginctx.FullPath() // form of /application/:app_id
		act := ginctx.Request.Method
		ok, err := authz.Enforce(acc.Sub, res.Workspace.Id, obj, act)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you have no permission to perform the action"})
			return
		}

		ctx = context.WithValue(ctx, authorizator.CtxWs, res.Workspace)
		ctx = context.WithValue(ctx, authorizator.CtxWst, res.WorkspaceTier)
		ginctx.Next()
	}
}
