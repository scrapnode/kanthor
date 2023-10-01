package sdkapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/debugging"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/sdkapi/docs"
	"github.com/scrapnode/kanthor/services/sdkapi/middlewares"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	idempotency idempotency.Idempotency,
	coordinator coordinator.Coordinator,
	metrics metric.Metrics,
	authz authorizator.Authorizator,
	uc usecase.Sdk,
) services.Service {
	logger = logger.With("service", "sdkapi")
	return &sdkapi{
		conf:        conf,
		logger:      logger,
		idempotency: idempotency,
		coordinator: coordinator,
		metrics:     metrics,
		authz:       authz,
		uc:          uc,

		debugger: debugging.NewServer(),
	}
}

type sdkapi struct {
	conf        *config.Config
	logger      logging.Logger
	idempotency idempotency.Idempotency
	coordinator coordinator.Coordinator
	metrics     metric.Metrics
	authz       authorizator.Authorizator
	uc          usecase.Sdk

	debugger debugging.Server
	server   *http.Server
}

func (service *sdkapi) Start(ctx context.Context) error {
	if err := service.debugger.Start(ctx); err != nil {
		return err
	}

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
		Addr:    service.conf.SdkApi.Gateway.Httpx.Addr,
		Handler: service.router(),
	}

	service.logger.Info("started")
	return nil
}

func (service *sdkapi) router() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	// system routes
	RegisterHealthcheck(router, service)

	swagger := router.Group("/swagger")
	{
		swagger.GET("/*any", ginswagger.WrapHandler(
			swaggerfiles.Handler,
			ginswagger.PersistAuthorization(true),
			ginswagger.InstanceName(docs.SwaggerInfoSdk.InfoInstanceName),
		))
	}

	api := router.Group("/api")
	{
		api.Use(ginmw.UseStartup())
		api.Use(ginmw.UseMetrics(service.metrics))
		api.Use(ginmw.UseIdempotency(service.logger, service.idempotency))
		api.Use(middlewares.UseAuthx(service.conf.SdkApi, service.logger, service.authz, service.uc))
		api.Use(ginmw.UsePaging(service.logger, 5, 30))

		RegisterAccountRoutes(api.Group("/account"), service)
		RegisterApplicationRoutes(api.Group("/application"), service)
		RegisterEndpointRoutes(api.Group("/application/:app_id/endpoint"), service)
		RegisterEndpointRuleRoutes(api.Group("/application/:app_id/endpoint/:ep_id/rule"), service)
		RegisterMessageRoutes(api.Group("/application/:app_id/message"), service)
	}

	return router
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

	if err := service.authz.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.metrics.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.debugger.Stop(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *sdkapi) Run(ctx context.Context) error {
	if err := service.coordinate(); err != nil {
		return err
	}

	go func() {
		if err := service.debugger.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			service.logger.Error(err)
		}
	}()

	service.logger.Infow("running", "addr", service.conf.SdkApi.Gateway.Httpx.Addr)
	if err := service.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		service.logger.Error(err)
	}

	return nil
}
