package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/sdkapi/docs"
	"github.com/scrapnode/kanthor/services/sdkapi/middlewares"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
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
	metrics metrics.Metrics,
	authz authorizator.Authorizator,
	uc usecase.Sdk,
) services.Service {
	logger = logger.With("service", "sdkapi")
	return &sdkapi{
		conf:        conf,
		logger:      logger,
		validator:   validator,
		idempotency: idempotency,
		coordinator: coordinator,
		metrics:     metrics,
		authz:       authz,
		uc:          uc,
	}
}

type sdkapi struct {
	conf        *config.Config
	logger      logging.Logger
	validator   validator.Validator
	idempotency idempotency.Idempotency
	coordinator coordinator.Coordinator
	metrics     metrics.Metrics
	authz       authorizator.Authorizator
	uc          usecase.Sdk

	server *http.Server
}

func (service *sdkapi) Start(ctx context.Context) error {
	if err := service.authz.Connect(ctx); err != nil {
		return err
	}

	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	if err := service.idempotency.Connect(ctx); err != nil {
		return err
	}

	if err := service.coordinator.Connect(ctx); err != nil {
		return err
	}

	service.build()

	service.logger.Info("started")
	return nil
}

func (service *sdkapi) build() {
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
	router.GET("/metrics", gin.WrapH(service.metrics.Handler()))

	swagger := router.Group("/swagger")
	{
		swagger.GET("/*any", ginswagger.WrapHandler(
			swaggerfiles.Handler,
			ginswagger.PersistAuthorization(true),
			ginswagger.InstanceName(docs.SwaggerInfoSdk.InfoInstanceName),
		))
	}
	// api routes
	api := router.Group("/api")
	{
		api.Use(ginmw.UseStartup())
		api.Use(ginmw.UseMetrics(service.metrics))
		api.Use(ginmw.UseIdempotency(service.logger, service.idempotency))
		api.Use(middlewares.UseAuthx(service.conf.SdkApi, service.logger, service.validator, service.authz, service.uc))
		api.Use(ginmw.UsePaging(service.logger, 5, 30))

		UseApplicationRoutes(api.Group("/application"), service)
		UseEndpointRoutes(api.Group("/application/:app_id/endpoint"), service)
		UseEndpointRuleRoutes(api.Group("/application/:app_id/endpoint/:ep_id/rule"), service)
		UseMessageRoutes(api.Group("/application/:app_id/message"), service)
	}

	service.server = &http.Server{
		Addr:    service.conf.SdkApi.Gateway.Httpx.Addr,
		Handler: router,
	}
}

func (service *sdkapi) Stop(ctx context.Context) error {
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

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.authz.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *sdkapi) Run(ctx context.Context) error {
	if err := service.coordinate(); err != nil {
		return err
	}

	service.logger.Infow("running", "addr", service.conf.SdkApi.Gateway.Httpx.Addr)

	err := service.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
