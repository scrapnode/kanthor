package sdkapi

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
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
	infra *infrastructure.Infrastructure,
	uc usecase.Sdk,
) services.Service {
	logger = logger.With("service", "sdkapi")
	return &sdkapi{
		conf:   conf,
		logger: logger,
		infra:  infra,
		uc:     uc,
	}
}

type sdkapi struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	uc     usecase.Sdk

	server *http.Server

	mu     sync.Mutex
	status int
}

func (service *sdkapi) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	service.server = &http.Server{
		Addr:    service.conf.SdkApi.Gateway.Addr,
		Handler: service.router(),
	}

	service.status = patterns.StatusStarted
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
		api.Use(ginmw.UseMetric(config.SERVICE_SDK_API, service.infra.Metric))
		api.Use(ginmw.UseIdempotency(service.logger, service.infra.Idempotency))
		api.Use(middlewares.UseAuth(service.infra.Authorizator, service.uc))
		api.Use(ginmw.UseAuthz(service.infra.Authorizator))
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
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status != patterns.StatusStarted {
		return ErrNotStarted
	}
	service.status = patterns.StatusStopped
	service.logger.Info("stopped")

	var returning error
	if err := service.server.Shutdown(ctx); err != nil {
		service.logger.Error(err)
		returning = errors.Join(returning, err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
		returning = errors.Join(returning, err)
	}

	if err := service.infra.Disconnect(ctx); err != nil {
		service.logger.Error(err)
		returning = errors.Join(returning, err)
	}

	return returning
}

func (service *sdkapi) Run(ctx context.Context) error {
	service.logger.Infow("running", "addr", service.conf.SdkApi.Gateway.Addr)
	if err := service.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		service.logger.Error(err)
	}

	return nil
}
