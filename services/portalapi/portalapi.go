package portalapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	httpxmw "github.com/scrapnode/kanthor/infrastructure/gateway/httpx/middlewares"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/portalapi/docs"
	"github.com/scrapnode/kanthor/services/portalapi/middlewares"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	validator validator.Validator,
	idempotency idempotency.Idempotency,
	auth authenticator.Authenticator,
	authz authorizator.Authorizator,
	uc usecase.Portal,
) services.Service {
	logger = logger.With("service", "portalapi")
	return &portalapi{
		conf:        conf,
		logger:      logger,
		validator:   validator,
		idempotency: idempotency,
		auth:        auth,
		authz:       authz,
		uc:          uc,
	}
}

type portalapi struct {
	conf        *config.Config
	logger      logging.Logger
	validator   validator.Validator
	idempotency idempotency.Idempotency
	auth        authenticator.Authenticator
	authz       authorizator.Authorizator
	uc          usecase.Portal

	server *http.Server
}

func (service *portalapi) Start(ctx context.Context) error {
	if err := service.authz.Connect(ctx); err != nil {
		return err
	}

	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	if err := service.idempotency.Connect(ctx); err != nil {
		return err
	}

	router := gin.New()
	router.Use(gin.Recovery())
	// system routes
	router.GET("/", func(ginctx *gin.Context) {
		ginctx.JSON(http.StatusOK, gin.H{"version": service.conf.Version})
	})
	router.GET("/readiness", func(ginctx *gin.Context) {
		// @TODO: add starting up checking here
		ginctx.String(http.StatusOK, "ready")
	})
	router.GET("/liveness", func(ginctx *gin.Context) {
		ginctx.String(http.StatusOK, "live")
	})
	swagger := router.Group("/swagger")
	{
		swagger.GET("/*any", ginswagger.WrapHandler(
			swaggerfiles.Handler,
			ginswagger.PersistAuthorization(true),
			ginswagger.InstanceName(docs.SwaggerInfoPortal.InfoInstanceName),
		))
	}
	// api routes
	api := router.Group("/api")
	{
		api.Use(httpxmw.UseStartup())
		api.Use(httpxmw.UseIdempotency(service.logger, service.idempotency))
		api.Use(middlewares.UseAuth(service.auth))
		api.Use(middlewares.UseAuthz(service.validator, service.authz, service.uc))
		api.Use(httpxmw.UsePaging(service.logger, 5, 30))

		UseWorkspaceRoutes(
			api.Group("/workspace"),
			service.logger, service.validator, service.uc,
		)

		UseWorkspaceCredentialsRoutes(
			api.Group("/workspace/me/credentials"),
			service.logger, service.validator, service.uc,
		)
	}

	service.server = &http.Server{
		Addr:    service.conf.PortalApi.Gateway.Httpx.Addr,
		Handler: router,
	}

	service.logger.Info("started")
	return nil
}

func (service *portalapi) Stop(ctx context.Context) error {
	service.logger.Info("stopped")

	if err := service.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := service.idempotency.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.authz.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *portalapi) Run(ctx context.Context) error {
	service.logger.Infow("running", "addr", service.conf.PortalApi.Gateway.Httpx.Addr)

	err := service.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
