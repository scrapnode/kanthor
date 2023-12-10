package rest

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/openapi"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/usecase"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	uc usecase.Portal,
) patterns.Runnable {
	logger = logger.With("service", "portal")
	return &portal{
		conf:   conf,
		logger: logger,
		infra:  infra,
		db:     db,
		uc:     uc,
	}
}

type portal struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	db     database.Database
	uc     usecase.Portal

	server *http.Server

	mu     sync.Mutex
	status int
}

func (service *portal) Start(ctx context.Context) error {
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
	// register authenticator
	for _, auth := range service.conf.Authenticator {
		if auth.Engine == authenticator.EngineAsk {
			service.infra.Authenticator.Register(auth.Engine, authenticator.NewAsk(auth.Ask))
		}
		if auth.Engine == authenticator.EngineForward {
			service.infra.Authenticator.Register(auth.Engine, authenticator.NewForward(auth.Forward))
		}
	}

	router, err := service.router()
	if err != nil {
		return err
	}
	service.server = &http.Server{
		Addr:    service.conf.Gateway.Addr,
		Handler: router,
	}

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *portal) router() (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: false,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Authorization-Engine", "X-Authorization-Workspace"},
		MaxAge:           time.Hour * 12,
	}))
	// system routes
	RegisterHealthcheck(router, service)

	swagger := router.Group("/swagger")
	{
		swagger.GET("/*any", ginswagger.WrapHandler(
			swaggerfiles.Handler,
			ginswagger.PersistAuthorization(true),
			ginswagger.InstanceName(openapi.SwaggerInfoPortal.InfoInstanceName),
		))
	}

	api := router.Group("/api")
	{
		api.Use(middlewares.UseStartup(&service.conf.Gateway))
		api.Use(middlewares.UseMetric(service.infra.Metric, "portal"))
		api.Use(middlewares.UseIdempotency(service.logger, service.infra.Idempotency))
		api.Use(middlewares.UsePaging(service.logger, 5, 30))

		api.Use(middlewares.UseAuth(service.infra.Authenticator, service.conf.Authenticator[0].Engine))

		// IMPORTANT: always put the longer route in the top
		RegisterAccountRoutes(api.Group("/account"), service)
		RegisterWorkspaceRoutes(api.Group("/workspace"), service)
	}

	return router, nil
}

func (service *portal) Stop(ctx context.Context) error {
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

func (service *portal) Run(ctx context.Context) error {
	service.logger.Infow("running", "addr", service.conf.Gateway.Addr)
	if err := service.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
