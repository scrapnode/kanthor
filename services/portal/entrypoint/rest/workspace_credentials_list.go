package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceCredentialsListRes struct {
	Data  []WorkspaceCredentials `json:"data"`
	Count int64                  `json:"count"`
}

// UseWorkspaceCredentialsList
// @Tags		workspace
// @Router		/workspace/me/credentials	[get]
// @Param		_q							query		string						false	"search keyword" 			minlength(2)  maxlength(32)
// @Param		_limit						query		int							false	"limit returning records"	minimum(5)    maximum(30)
// @Param		_page						query		int							false	"requesting page"
// @Success		200							{object}	WorkspaceCredentialsListRes
// @Failure		default						{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceCredentialsList(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		var query gateway.Query
		if err := ginctx.BindQuery(&query); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("unable to parse your request query"))
			return
		}

		in := &usecase.WorkspaceCredentialsListIn{
			WsId:   ws.Id,
			Search: query.Search,
			Limit:  query.Limit,
			Page:   query.Page,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.WorkspaceCredentials().List(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceCredentialsListRes{Data: make([]WorkspaceCredentials, 0), Count: out.Count}
		for _, ws := range out.Data {
			res.Data = append(res.Data, *ToWorkspaceCredentials(&ws))
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
