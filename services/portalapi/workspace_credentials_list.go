package portalapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validation"
	"github.com/scrapnode/kanthor/pkg/utils"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
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
func UseWorkspaceCredentialsList(logger logging.Logger, validator validation.Validator, uc usecase.Portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		ucreq := &usecase.WorkspaceCredentialsListReq{ListReq: ginctx.MustGet("list_req").(*structure.ListReq)}
		if err := validator.Struct(ucreq); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.WorkspaceCredentials().List(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceCredentialsListRes{ListRes: ucres.ListRes}
		ginctx.JSON(http.StatusOK, res)
	}
}
