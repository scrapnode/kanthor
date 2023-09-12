package sdkapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/pkg/utils"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

type EndpointRuleDeleteRes struct {
	*entities.EndpointRule
}

// UseEndpointRuleDelete
// @Tags		endpoint rule
// @Router		/application/{app_id}/endpoint/{ep_id}/rule/{epr_id}	[delete]
// @Param		app_id													path		string					true	"application id"
// @Param		ep_id													path		string					true	"endpoint id"
// @Param		epr_id													path		string					true	"rule id"
// @Success		200														{object}	EndpointRuleDeleteRes
// @Failure		default													{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointRuleDelete(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		appId := ginctx.Param("app_id")
		epId := ginctx.Param("ep_id")
		id := ginctx.Param("epr_id")
		ucreq := &usecase.EndpointRuleDeleteReq{AppId: appId, EpId: epId, Id: id}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.EndpointRule().Delete(ctx, ucreq)
		if err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointRuleDeleteRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
