package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceListRes struct {
	Data []entities.Workspace `json:"data"`
}

// UseWorkspaceList
// @Tags		workspace
// @Router		/workspace				[get]
// @Success		200						{object}	WorkspaceGetRes
// @Failure		default					{object}	gateway.Error
// @Security	BearerAuth
func UseWorkspaceList(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		in := &usecase.WorkspaceListIn{AccId: acc.Sub}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Workspace().List(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceListRes{Data: []entities.Workspace{}}
		for _, ws := range out.Workspaces {
			res.Data = append(res.Data, *ws)
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
