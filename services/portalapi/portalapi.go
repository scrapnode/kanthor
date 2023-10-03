package portalapi

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/debugging"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
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
		idempotency: idempotency,
		coordinator: coordinator,
		metrics:     metrics,
		auth:        auth,
		authz:       authz,
		uc:          uc,
		debugger:    debugging.NewServer(),
	}
}

type portalapi struct {
	conf        *config.Config
	logger      logging.Logger
	idempotency idempotency.Idempotency
	coordinator coordinator.Coordinator
	metrics     metric.Metrics
	auth        authenticator.Authenticator
	authz       authorizator.Authorizator
	uc          usecase.Portal

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
		api.Use(ginmw.UseMetrics(service.metrics))
		api.Use(ginmw.UseIdempotency(service.logger, service.idempotency))
		api.Use(ginmw.UsePaging(service.logger, 5, 30))
		api.Use(middlewares.UseAuth(service.auth))

		RegisterAccountRoutes(api.Group("/account"), service)
		RegisterWorkspaceRoutes(api.Group("/workspace").Use(middlewares.UseAuthz(service.authz, service.uc)), service)
		RegisterWorkspaceCredentialsRoutes(api.Group("/workspace/me/credentials").Use(middlewares.UseAuthz(service.authz, service.uc)), service)
	}

	return router
}

func (service *portalapi) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

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

func (service *portalapi) Run(ctx context.Context) error {
	if err := service.coordinate(); err != nil {
		return err
	}

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
