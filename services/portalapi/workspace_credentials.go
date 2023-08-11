package portalapi

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
)

func UseWorkspaceCredentialsRoutes(
	router *gin.RouterGroup,
	logger logging.Logger,
	validator validator.Validator,
	uc usecase.Portal,
) {
	router.GET("", UseWorkspaceCredentialsList(logger, validator, uc))
	router.POST("", UseWorkspaceCredentialsCreate(logger, validator, uc))
	router.GET("/:wsc_id", UseWorkspaceCredentialsGet(logger, validator, uc))
	router.PUT("/:wsc_id", UseWorkspaceCredentialsUpdate(logger, validator, uc))
	router.PUT("/:wsc_id/expiration", UseWorkspaceCredentialsExpire(logger, validator, uc))
}
