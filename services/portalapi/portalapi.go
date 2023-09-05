package portalapi

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
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
	coordinator coordinator.Coordinator,
	metrics metric.Metrics,
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
		coordinator: coordinator,
		metrics:     metrics,
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
	coordinator coordinator.Coordinator
	metrics     metric.Metrics
	auth        authenticator.Authenticator
	authz       authorizator.Authorizator
	uc          usecase.Portal

	server *http.Server
}

func (service *portalapi) Start(ctx context.Context) error {
	if err := service.metrics.Connect(ctx); err != nil {
		return err
	}

	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	if err := service.authz.Connect(ctx); err != nil {
		return err
	}

	if err := service.idempotency.Connect(ctx); err != nil {
		return err
	}

	if err := service.coordinator.Connect(ctx); err != nil {
		return err
	}

	service.server = &http.Server{
		Addr:    service.conf.PortalApi.Gateway.Httpx.Addr,
		Handler: service.router(),
	}

	service.logger.Info("started")
	return nil
}

func (service *portalapi) router() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())
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
		api.Use(ginmw.UseStartup())
		api.Use(ginmw.UseMetrics(service.metrics))
		api.Use(ginmw.UseIdempotency(service.logger, service.idempotency))
		api.Use(ginmw.UsePaging(service.logger, 5, 30))
		api.Use(middlewares.UseAuth(service.auth))

		UseAccountRoutes(api.Group("/account"), service)
		UseWorkspaceRoutes(api.Group("/workspace").Use(middlewares.UseAuthz(service.validator, service.authz, service.uc)), service)
		UseWorkspaceCredentialsRoutes(api.Group("/workspace/me/credentials").Use(middlewares.UseAuthz(service.validator, service.authz, service.uc)), service)
	}

	return router
}

func (service *portalapi) Stop(ctx context.Context) error {
	service.logger.Info("stopped")

	if err := service.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := service.coordinator.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.idempotency.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.authz.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.metrics.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *portalapi) Run(ctx context.Context) error {
	if err := service.coordinate(); err != nil {
		return err
	}

	service.logger.Infow("running", "addr", service.conf.PortalApi.Gateway.Httpx.Addr)
	if err := service.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
