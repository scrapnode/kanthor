package rest

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/authenticator"
	"github.com/scrapnode/kanthor/database"
	ginmw "github.com/scrapnode/kanthor/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/entrypoint/rest/docs"
	"github.com/scrapnode/kanthor/services/portal/entrypoint/rest/middlewares"
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
	router.Use(cors.Default())
	// system routes
	RegisterHealthcheck(router, service)

	swagger := router.Group("/swagger")
	{
		swagger.GET("/*any", ginswagger.WrapHandler(
			swaggerfiles.Handler,
			ginswagger.PersistAuthorization(true),
			ginswagger.InstanceName(docs.SwaggerInfoPortal.InfoInstanceName),
		))
	}

	api := router.Group("/api")
	{
		api.Use(ginmw.UseStartup(&service.conf.Gateway))
		api.Use(ginmw.UseMetric(service.infra.Metric, "portal"))
		api.Use(ginmw.UseIdempotency(service.logger, service.infra.Idempotency))
		api.Use(ginmw.UsePaging(service.logger, 5, 30))
		auth, err := authenticator.New(&service.conf.Authenticator, service.logger)
		if err != nil {
			return nil, err
		}
		api.Use(middlewares.UseAuth(auth, service.uc))

		RegisterAccountRoutes(api.Group("/account"), service)
		RegisterWorkspaceRoutes(api.Group("/workspace"), service)
		RegisterWorkspaceCredentialsRoutes(api.Group("/workspace/me/credentials"), service)
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