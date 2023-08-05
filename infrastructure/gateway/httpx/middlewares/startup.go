package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func UseStartup() gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		// hard limit with 30s timeout
		// if the API need more than 30 seconds to execute, shut the f*ck up and refactor it
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		ginctx.Set("ctx", ctx)
		ginctx.Next()
	}
}
