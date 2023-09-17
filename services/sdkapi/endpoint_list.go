package sdkapi

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
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
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
func UseEndpointList(logger logging.Logger, validator validation.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		appId := ginctx.Param("app_id")

		ucreq := &usecase.EndpointListReq{
			AppId:   appId,
			ListReq: ginctx.MustGet("list_req").(*structure.ListReq),
		}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Endpoint().List(ctx, ucreq)
		if err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointListRes{ListRes: ucres.ListRes}
		ginctx.JSON(http.StatusOK, res)
	}
}
