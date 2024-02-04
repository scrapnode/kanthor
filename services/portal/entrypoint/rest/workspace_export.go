package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceExportRes struct {
	*WorkspaceSnapshot
} // @name WorkspaceExportRes

// UseWorkspaceExport
// @Tags			workspace
// @Router		/workspace/{ws_id}/transfer	[get]
// @Param			ws_id												path			string						true	"workspace id"
// @Success		200													{object}	WorkspaceExportRes
// @Failure		default											{object}	gateway.Err
// @Security	Authorization
func UseWorkspaceExport(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		getin := &usecase.WorkspaceGetIn{
			AccId: acc.Sub,
			Id:    ginctx.Param("ws_id"),
		}
		if err := getin.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}
		getout, err := service.uc.Workspace().Get(ctx, getin)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		exportin := &usecase.WorkspaceExportIn{Id: getout.Doc.Id}
		if err := exportin.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}
		exportout, err := service.uc.Workspace().Export(ctx, exportin)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &WorkspaceExportRes{ToWorkspaceSnapshot(exportout.Data)}
		ginctx.JSON(http.StatusOK, res)
	}
}
