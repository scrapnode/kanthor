package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceImportReq struct {
	Snapshot *WorkspaceSnapshot `json:"snapshot" form:"snapshot" binding:"required"`
} // @name WorkspaceImportReq

type WorkspaceImportRes struct {
	AppIds []string `json:"app_id"`
	EpIds  []string `json:"ep_id"`
	EprIds []string `json:"epr_id"`
} // @name WorkspaceImportRes

// UseWorkspaceImport
// @Tags		workspace
// @Router		/workspace/{ws_id}/transfer	[post]
// @Param		ws_id						path		string						true	"workspace id"
// @Param		payload						body		WorkspaceImportReq	true	"import payload"
// @Success		200							{object}	WorkspaceImportRes
// @Failure		default						{object}	gateway.Err
// @Security	Authorization
func UseWorkspaceImport(service *portal) gin.HandlerFunc {
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

		var req WorkspaceImportReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}
		importin := &usecase.WorkspaceImportIn{
			Id:       getout.Doc.Id,
			Snapshot: FromWorkspaceSnapshot(req.Snapshot, getout.Doc.Id),
		}
		if err := importin.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}
		importout, err := service.uc.Workspace().Import(ctx, importin)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &WorkspaceImportRes{
			AppIds: importout.AppIds,
			EpIds:  importout.EpIds,
			EprIds: importout.EprIds,
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
