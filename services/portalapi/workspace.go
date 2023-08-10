package portalapi

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
)

func UseWorkspaceRoutes(router *gin.RouterGroup, logger logging.Logger, validator validator.Validator, uc usecase.Portal) {
	router.GET("/me", UseWorkspaceGet())
}
