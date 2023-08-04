package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

type endpointRuleListRes struct {
	*structure.ListRes[entities.EndpointRule]
}

func UseEndpointRuleList(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet("ctx").(context.Context)
		appId := ginctx.Param("app_id")
		epId := ginctx.Param("ep_id")

		ucreq := &usecase.EndpointRuleListReq{
			AppId:   appId,
			EpId:    epId,
			ListReq: ginctx.MustGet("list_req").(*structure.ListReq),
		}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		ucres, err := uc.EndpointRule().List(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops, something went wrong"})
			return
		}

		res := &endpointRuleListRes{ListRes: ucres.ListRes}
		ginctx.JSON(http.StatusOK, res)
	}
}
