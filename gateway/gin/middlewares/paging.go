package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
)

type paging struct {
	Cursor string `form:"_cursor" binding:"omitempty,min=29,max=32"`
	Search string `form:"_q" binding:"omitempty,min=2,max=32"`
	Limit  int    `form:"_limit" binding:"omitempty,min=5,max=30"`

	Ids []string `form:"_id"`
}

func UsePaging(logger logging.Logger, minLimit, maxLimit int) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		if ginctx.Request.Method != http.MethodGet {
			ginctx.Next()
			return
		}

		var p paging
		if err := ginctx.BindQuery(&p); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("unable to parse your request query"))
			return
		}
		p.Limit = utils.MinInt(utils.MaxInt(minLimit, p.Limit), maxLimit)

		ginctx.Set("list_req", &structure.ListReq{Cursor: p.Cursor, Search: p.Search, Limit: p.Limit, Ids: p.Ids})
		ginctx.Next()
	}
}
