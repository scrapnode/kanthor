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

type endpointRuleCreateReq struct {
	Name string `json:"name" binding:"required"`

	Priority            int32  `json:"priority"`
	Exclusionary        bool   `json:"exclusionary"`
	ConditionSource     string `json:"condition_source" binding:"required"`
	ConditionExpression string `json:"condition_expression" binding:"required"`
}

type endpointRuleCreateRes struct {
	*entities.EndpointRule
}

func UseEndpointRuleCreate(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req endpointRuleCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
			return
		}

		ctx := ginctx.MustGet("ctx").(context.Context)
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
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		ucres, err := uc.EndpointRule().Create(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops, something went wrong"})
			return
		}

		res := &endpointRuleCreateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
