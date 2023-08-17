package portalapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/services/command"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
	"net/http"
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
	validator validator.Validator,
	uc portaluc.Portal,
	coord coordinator.Coordinator,
) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCredentialsCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

		ucreq := &portaluc.WorkspaceCredentialsGenerateReq{WorkspaceId: ws.Id, Count: 1}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.WorkspaceCredentials().Generate(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		err = coord.Send(
			ctx,
			command.WorkspaceCredentialsCreated,
			&command.WorkspaceCredentialsCreatedReq{Docs: ucres.Credentials},
		)
		if err != nil {
			logger.Error(err)
		}

		res := &WorkspaceCredentialsCreateRes{
			Id:       ucres.Credentials[0].Id,
			User:     ucres.Credentials[0].Id,
			Password: ucres.Passwords[ucres.Credentials[0].Id],
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
