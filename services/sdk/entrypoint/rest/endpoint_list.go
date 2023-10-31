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

type EndpointListRes struct {
	*structure.ListRes[entities.Endpoint]
}

// UseEndpointList
// @Tags		endpoint
// @Router		/application/{app_id}/endpoint	[get]
// @Param		app_id							path		string					true	"application id"
// @Param		_cursor							query		string					false	"current query cursor"					minlength(29) maxlength(32)
// @Param		_q								query		string					false	"search keyword" 						minlength(2)  maxlength(32)
// @Param		_limit							query		int						false	"limit returning records"				minimum(5)    maximum(30)
// @Param		_id								query		[]string				false	"only return records with selected ids"
// @Success		200								{object}	EndpointListRes
// @Failure		default							{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointList(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		appId := ginctx.Param("app_id")

		ucreq := &usecase.EndpointListReq{
			AppId:   appId,
			ListReq: ginctx.MustGet("list_req").(*structure.ListReq),
		}
		if err := ucreq.Validate(); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := service.uc.Endpoint().List(ctx, ucreq)
		if err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointListRes{ListRes: ucres.ListRes}
		ginctx.JSON(http.StatusOK, res)
	}
}