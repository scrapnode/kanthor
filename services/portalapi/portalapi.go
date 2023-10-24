package portalapi

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/database"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/portalapi/docs"
	"github.com/scrapnode/kanthor/services/portalapi/middlewares"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	auth authenticator.Authenticator,
	infra *infrastructure.Infrastructure,
	db database.Database,
	uc usecase.Portal,
) services.Service {
	logger = logger.With("service", "portalapi")
	return &portalapi{
		conf:   conf,
		logger: logger,
		auth:   auth,
		infra:  infra,
		db:     db,
		uc:     uc,
	}
}

type portalapi struct {
	conf   *config.Config
	logger logging.Logger
	auth   authenticator.Authenticator
	infra  *infrastructure.Infrastructure
	db     database.Database
	uc     usecase.Portal

	server *http.Server

	mu     sync.Mutex
	status int
}

func (service *portalapi) Start(ctx context.Context) error {
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
		Addr:    service.conf.PortalApi.Gateway.Addr,
		Handler: service.router(),
	}

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *portalapi) router() *gin.Engine {
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
		api.Use(ginmw.UseStartup(&service.conf.PortalApi.Gateway))
		api.Use(ginmw.UseMetric(config.SERVICE_PORTAL_API, service.infra.Metric))
		api.Use(ginmw.UseIdempotency(service.logger, service.infra.Idempotency))
		api.Use(ginmw.UsePaging(service.logger, 5, 30))
		api.Use(middlewares.UseAuth(service.auth, service.uc))

		RegisterAccountRoutes(api.Group("/account"), service)
		RegisterWorkspaceRoutes(api.Group("/workspace"), service)
		RegisterWorkspaceCredentialsRoutes(api.Group("/workspace/me/credentials"), service)
	}

	return router
}

func (service *portalapi) Stop(ctx context.Context) error {
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

func (service *portalapi) Run(ctx context.Context) error {
	service.logger.Infow("running", "addr", service.conf.PortalApi.Gateway.Addr)
	if err := service.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
