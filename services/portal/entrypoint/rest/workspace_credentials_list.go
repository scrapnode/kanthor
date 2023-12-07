package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/structure"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceCredentialsListRes struct {
	*structure.ListRes[entities.WorkspaceCredentials]
}

// UseWorkspaceCredentialsList
// @Tags		workspace
// @Router		/workspace/me/credentials	[get]
// @Param		_cursor						query		string					false	"current query cursor"					minlength(29) maxlength(32)
// @Param		_q							query		string					false	"search keyword" 						minlength(2)  maxlength(32)
// @Param		_limit						query		int						false	"limit returning records"				minimum(5)    maximum(30)
// @Param		_id							query		[]string				false	"only return records with selected ids"
// @Success		200							{object}	WorkspaceCredentialsListRes
// @Failure		default						{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceCredentialsList(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)

		in := &usecase.WorkspaceCredentialsListIn{ListReq: ginctx.MustGet("list_req").(*structure.ListReq)}
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

		res := &WorkspaceCredentialsListRes{ListRes: out.ListRes}
		ginctx.JSON(http.StatusOK, res)
	}
}
