package rest

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/openapi"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/project"
	"github.com/scrapnode/kanthor/services/sdk/config"
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
	// register the SDK internal authenticator
	err := service.infra.Authenticator.Register(AuthzEngineInternal, &internal{uc: service.uc})
	if err != nil {
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
	router.Use(middlewares.UseCors(service.conf.Gateway.Origins))
	// system routes
	RegisterHealthcheck(router, service)

	swagger := router.Group("/swagger")
	{
		swagger.GET("/*any", ginswagger.WrapHandler(
			swaggerfiles.Handler,
			ginswagger.PersistAuthorization(true),
			ginswagger.InstanceName(openapi.SwaggerInfoSdk.InfoInstanceName),
		))
	}

	api := router.Group("/api")
	{
		api.Use(middlewares.UseStartup(&service.conf.Gateway))
		api.Use(middlewares.UseMetric(service.infra.Metric, "sdk"))
		api.Use(middlewares.UseIdempotency(service.logger, service.infra.Idempotency, project.IsDev()))

		api.Use(middlewares.UseAuth(service.infra.Authenticator, AuthzEngineInternal))

		// IMPORTANT: always put the longer route in the top
		RegisterAccountRoutes(api.Group("/account"), service)

		RegisterEndpointRuleRoutes(api.Group("/rule"), service)
		RegisterEndpointRoutes(api.Group("/endpoint"), service)
		RegisterMessageRoutes(api.Group("/message"), service)
		RegisterApplicationRoutes(api.Group("/application"), service)
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
