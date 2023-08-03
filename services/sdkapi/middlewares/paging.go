package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"net/http"
)

type paging struct {
	Cursor string `form:"_cursor"`
	Search string `form:"_q"`
	Limit  int    `form:"_limit"`

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
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to parse your request query"})
			return
		}
		p.Limit = utils.MinInt(utils.MaxInt(minLimit, p.Limit), maxLimit)

		ginctx.Set("list_req", &structure.ListReq{Cursor: p.Cursor, Search: p.Search, Limit: p.Limit, Ids: p.Ids})
		ginctx.Next()
	}
}
