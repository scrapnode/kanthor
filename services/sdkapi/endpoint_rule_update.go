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

type EndpointRuleUpdateReq struct {
	Name string `json:"name" binding:"required"`
}

type EndpointRuleUpdateRes struct {
	*entities.EndpointRule
}

// UseEndpointRuleUpdate
// @Tags		endpoint rule
// @Router		/application/{app_id}/endpoint/{ep_id}/rule/{epr_id}	[put]
// @Param		app_id													path		string					true	"application id"
// @Param		ep_id													path		string					true	"endpoint id"
// @Param		epr_id													path		string					true	"rule id"
// @Param		props													body		EndpointRuleUpdateReq	true	"rule properties"
// @Success		200														{object}	EndpointRuleUpdateRes
// @Failure		default													{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointRuleUpdate(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointRuleUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		appId := ginctx.Param("app_id")
		epId := ginctx.Param("ep_id")
		id := ginctx.Param("epr_id")
		ucreq := &usecase.EndpointRuleUpdateReq{AppId: appId, EpId: epId, Id: id, Name: req.Name}
		if err := validator.Struct(ucreq); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.EndpointRule().Update(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointRuleUpdateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
