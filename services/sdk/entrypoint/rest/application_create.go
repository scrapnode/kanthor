package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type ApplicationCreateReq struct {
	Name string `json:"name" binding:"required"`
}

type ApplicationCreateRes struct {
	*Application
}

// UseApplicationCreate
// @Tags		application
// @Router		/application		[post]
// @Param		props				body		ApplicationCreateReq	true	"application properties"
// @Success		201					{object}	ApplicationCreateRes
// @Failure		default				{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseApplicationCreate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req ApplicationCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.ApplicationCreateIn{
			WsId: ws.Id,
			Name: req.Name,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Application().Create(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &ApplicationCreateRes{ToApplication(out.Doc)}
		ginctx.JSON(http.StatusCreated, res)
	}
}
