package portalapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
	"net/http"
)

type WorkspaceUpdateReq struct {
	Name string `json:"name" binding:"required"`
}

type WorkspaceUpdateRes struct {
	*entities.Workspace
}

// UseWorkspaceUpdate
// @Tags		workspace
// @Router		/workspace/me			[put]
// @Param		props					body		WorkspaceUpdateReq	true	"credentials properties"
// @Success		200						{object}	WorkspaceUpdateRes
// @Failure		default					{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceUpdate(logger logging.Logger, validator validator.Validator, uc portaluc.Portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

		ucreq := &portaluc.WorkspaceUpdateReq{Id: ws.Id, Name: req.Name}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Workspace().Update(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceUpdateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}