package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

type endpointRuleDeleteRes struct {
	*entities.EndpointRule
}

func UseEndpointRuleDelete(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet("ctx").(context.Context)
		appId := ginctx.Param("app_id")
		epId := ginctx.Param("ep_id")
		id := ginctx.Param("epr_id")
		ucreq := &usecase.EndpointRuleDeleteReq{AppId: appId, EpId: epId, Id: id}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		ucres, err := uc.EndpointRule().Delete(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops, something went wrong"})
			return
		}

		res := &endpointRuleDeleteRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
