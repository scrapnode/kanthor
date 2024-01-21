package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/permissions"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceCredentialsCreateReq struct {
	Name      string `json:"name" default:"swagger demo"`
	ExpiredAt int64  `json:"expired_at" default:"1893456000000"`
} // @name WorkspaceCredentialsCreateReq

type WorkspaceCredentialsCreateRes struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
} // @name WorkspaceCredentialsCreateRes

// UseWorkspaceCredentialsCreate
// @Tags		credentials
// @Router		/credentials	[post]
// @Param		payload			body		WorkspaceCredentialsCreateReq	true	"credentials payload"
// @Success		200				{object}	WorkspaceCredentialsCreateRes
// @Failure		default			{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseWorkspaceCredentialsCreate(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCredentialsCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.WorkspaceCredentialsGenerateIn{
			WsId:      ws.Id,
			Name:      req.Name,
			ExpiredAt: req.ExpiredAt,
			// @TODO remove hardcode
			Role:        permissions.SdkOwner,
			Permissions: permissions.SdkOwnerPermissions,
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.WorkspaceCredentials().Generate(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		if err := service.infra.Authorizator.Refresh(ctx); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &WorkspaceCredentialsCreateRes{
			Id:       out.Credentials.Id,
			Name:     out.Credentials.Name,
			User:     out.Credentials.Id,
			Password: out.Password,
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
