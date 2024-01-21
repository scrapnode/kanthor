package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type ApplicationCreateReq struct {
	Name string `json:"name"`
} // @name ApplicationCreateReq

type ApplicationCreateRes struct {
	*Application
} // @name ApplicationCreateRes

// UseApplicationCreate
// @Tags		application
// @Router		/application		[post]
// @Param		payload				body		ApplicationCreateReq	true	"application payload"
// @Success		201					{object}	ApplicationCreateRes
// @Failure		default				{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseApplicationCreate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req ApplicationCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.ApplicationCreateIn{
			WsId: ws.Id,
			Name: req.Name,
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Application().Create(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &ApplicationCreateRes{ToApplication(out.Doc)}
		ginctx.JSON(http.StatusCreated, res)
	}
}
