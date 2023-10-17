package portalapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
)

type WorkspaceCredentialsCreateReq struct {
	Name      string `json:"name" binding:"required"`
	ExpiredAt int64  `json:"expired_at"`
}

type WorkspaceCredentialsCreateRes struct {
	Id       string `json:"id"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// UseWorkspaceCredentialsCreate
// @Tags		workspace
// @Router		/workspace/me/credentials	[post]
// @Param		props						body		WorkspaceCredentialsCreateReq	true	"credentials properties"
// @Success		200							{object}	WorkspaceCredentialsCreateRes
// @Failure		default						{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceCredentialsCreate(
	logger logging.Logger,
	uc portaluc.Portal,
) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCredentialsCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)

		ucreq := &portaluc.WorkspaceCredentialsGenerateReq{WorkspaceId: ws.Id, Name: req.Name, ExpiredAt: req.ExpiredAt}
		if err := ucreq.Validate(); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.WorkspaceCredentials().Generate(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceCredentialsCreateRes{
			Id:       ucres.Credentials.Id,
			User:     ucres.Credentials.Id,
			Password: ucres.Password,
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
