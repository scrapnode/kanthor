package rest

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/database"
	ginmw "github.com/scrapnode/kanthor/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/entrypoint/rest/docs"
	"github.com/scrapnode/kanthor/services/sdk/entrypoint/rest/middlewares"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	uc usecase.Sdk,
) patterns.Runnable {
	logger = logger.With("service", "sdk", "entrypoint", "rest")
	return &sdk{
		conf:   conf,
		logger: logger,
		infra:  infra,
		db:     db,
		uc:     uc,
	}
}

type sdk struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	db     database.Database
	uc     usecase.Sdk

	server *http.Server

	mu     sync.Mutex
	status int
}

func (service *sdk) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	if err := service.db.Connect(ctx); err != nil {
		return err
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	service.server = &http.Server{
		Addr:    service.conf.Gateway.Addr,
		Handler: service.router(),
	}

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *sdk) router() *gin.Engine {
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
		api.Use(ginmw.UseStartup(&service.conf.Gateway))
		api.Use(ginmw.UseMetric(service.infra.Metric, "sdk"))
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

func (service *sdk) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status != patterns.StatusStarted {
		return ErrNotStarted
	}
	service.status = patterns.StatusStopped
	service.logger.Info("stopped")

	var returning error
	if err := service.server.Shutdown(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.infra.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.db.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	return returning
}

func (service *sdk) Run(ctx context.Context) error {
	service.logger.Infow("running", "addr", service.conf.Gateway.Addr)
	if err := service.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
