package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func UseStartup() gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		ginctx.Set("ctx", ctx)
		ginctx.Next()
	}
}
