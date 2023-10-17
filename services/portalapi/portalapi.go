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
	"github.com/scrapnode/kanthor/infrastructure/debugging"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure/logging"
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
	infra *infrastructure.Infrastructure,
	auth authenticator.Authenticator,
	uc usecase.Portal,
) services.Service {
	logger = logger.With("service", "portalapi")
	return &portalapi{
		conf:     conf,
		logger:   logger,
		infra:    infra,
		auth:     auth,
		uc:       uc,
		debugger: debugging.NewServer(),
	}
}

type portalapi struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	auth   authenticator.Authenticator
	uc     usecase.Portal

	mu       sync.Mutex
	debugger debugging.Server
	server   *http.Server
}

func (service *portalapi) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if err := service.debugger.Start(ctx); err != nil {
		return err
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	if err := service.uc.Connect(ctx); err != nil {
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
		api.Use(ginmw.UseStartup())
		api.Use(ginmw.UseMetric(services.SERVICE_PORTAL_API, service.infra.Metric))
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

func (service *portalapi) Run(ctx context.Context) error {
	go func() {
		if err := service.debugger.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			service.logger.Error(err)
		}
	}()

	service.logger.Infow("running", "addr", service.conf.PortalApi.Gateway.Httpx.Addr)
	if err := service.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
