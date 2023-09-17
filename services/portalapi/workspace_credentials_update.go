package portalapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validation"
	"github.com/scrapnode/kanthor/pkg/utils"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
)

type WorkspaceCredentialsUpdateReq struct {
	Name string `json:"name"`
}

type WorkspaceCredentialsUpdateRes struct {
	*entities.WorkspaceCredentials
}

// UseWorkspaceCredentialsUpdate
// @Tags		workspace
// @Router		/workspace/me/credentials/{wsc_id}	[put]
// @Param		wsc_id								path		string							true	"credentials id"
// @Param		props								body		WorkspaceCredentialsUpdateReq	true	"credentials properties"
// @Success		200									{object}	WorkspaceCredentialsUpdateRes
// @Failure		default								{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceCredentialsUpdate(logger logging.Logger, validator validation.Validator, uc portaluc.Portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCredentialsUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

		id := ginctx.Param("wsc_id")
		ucreq := &portaluc.WorkspaceCredentialsUpdateReq{
			WorkspaceId: ws.Id,
			Id:          id,
			Name:        req.Name,
		}
		if err := validator.Struct(ucreq); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.WorkspaceCredentials().Update(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceCredentialsUpdateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
