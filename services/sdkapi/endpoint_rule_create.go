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

type EndpointRuleCreateReq struct {
	Name string `json:"name" binding:"required"`

	Priority            int32  `json:"priority"`
	Exclusionary        bool   `json:"exclusionary"`
	ConditionSource     string `json:"condition_source" binding:"required"`
	ConditionExpression string `json:"condition_expression" binding:"required"`
}

type EndpointRuleCreateRes struct {
	*entities.EndpointRule
}

// UseEndpointRuleCreate
// @Tags		endpoint rule
// @Router		/application/{app_id}/endpoint/{ep_id}/rule	[post]
// @Param		app_id										path		string					true	"application id"
// @Param		ep_id										path		string					true	"endpoint id"
// @Param		props										body		EndpointRuleCreateReq	true	"rule properties"
// @Success		201											{object}	EndpointRuleCreateRes
// @Failure		default										{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointRuleCreate(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointRuleCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		epId := ginctx.Param("ep_id")

		ucreq := &usecase.EndpointRuleCreateReq{
			EpId:                epId,
			Name:                req.Name,
			Priority:            req.Priority,
			Exclusionary:        req.Exclusionary,
			ConditionSource:     req.ConditionSource,
			ConditionExpression: req.ConditionExpression,
		}
		if err := validator.Struct(ucreq); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.EndpointRule().Create(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointRuleCreateRes{ucres.Doc}
		ginctx.JSON(http.StatusCreated, res)
	}
}
