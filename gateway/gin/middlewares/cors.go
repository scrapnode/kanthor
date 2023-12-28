package middlewares

import (
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func UseCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			if u.Hostname() == "localhost" {
				return true
			}

			return false
		},
		AllowCredentials: true,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"HEAD",
			"OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Idempotency-Key",
			"X-Authorization-Engine",
			"X-Authorization-Workspace",
		},
		MaxAge: time.Hour * 12,
	})
}
