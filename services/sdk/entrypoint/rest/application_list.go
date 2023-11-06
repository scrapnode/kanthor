package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type ApplicationListReq struct {
	*structure.ListRes[entities.Application]
}

// UseApplicationList
// @Tags		application
// @Router		/application		[get]
// @Param		_cursor				query		string					false	"current query cursor"					minlength(29) maxlength(32)
// @Param		_q					query		string					false	"search keyword" 						minlength(2)  maxlength(32)
// @Param		_limit				query		int						false	"limit returning records"				minimum(5)    maximum(30)
// @Param		_id					query		[]string				false	"only return records with selected ids"
// @Success		200					{object}	ApplicationListReq
// @Failure		default				{object}	gateway.Error
// @Security	BasicAuth
func UseApplicationList(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		in := &usecase.ApplicationListIn{ListReq: ginctx.MustGet("list_req").(*structure.ListReq)}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Application().List(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &ApplicationListReq{ListRes: out.ListRes}
		ginctx.JSON(http.StatusOK, res)
	}
}
