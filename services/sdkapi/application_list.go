package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

type ApplicationListReq struct {
	*structure.ListReq
}

type ApplicationListRes struct {
	*structure.ListRes[entities.Application]
}

func UseApplicationList(logger logging.Logger, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		req := &ApplicationListReq{ListReq: ginctx.MustGet("list_req").(*structure.ListReq)}

		ctx := ginctx.MustGet("ctx").(context.Context)
		ucreq := &usecase.ApplicationListReq{ListReq: req.ListReq}
		ucres, err := uc.Application().List(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops, something went wrong"})
			return
		}

		res := &ApplicationListRes{ListRes: ucres.ListRes}
		ginctx.JSON(http.StatusOK, res)
	}
}
