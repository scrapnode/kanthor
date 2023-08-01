package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/sdkapi/middlewares"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	authz authorizator.Authorizator,
	uc usecase.Sdk,
) services.Service {
	logger = logger.With("service", "sdkapi")
	return &sdkapi{
		conf:   conf,
		logger: logger,
		authz:  authz,
		uc:     uc,
	}
}

type sdkapi struct {
	conf   *config.Config
	logger logging.Logger
	authz  authorizator.Authorizator
	uc     usecase.Sdk

	server *http.Server
}

func (service *sdkapi) Start(ctx context.Context) error {
	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/readiness", func(ginctx *gin.Context) {
		// @TODO: add starting up checking here
		ginctx.String(http.StatusOK, "ready")
	})
	router.GET("/liveness", func(ginctx *gin.Context) {
		ginctx.String(http.StatusOK, "live")
	})
	router.Use(middlewares.UseAuth(service.uc, publicable))
	router.Use(middlewares.UseAuthz(service.authz, protectable))

	service.server = &http.Server{
		Addr:    service.conf.SdkApi.Gateway.Httpx.Addr,
		Handler: router,
	}

	service.logger.Info("started")
	return nil
}

func (service *sdkapi) Stop(ctx context.Context) error {
	service.logger.Info("stopped")

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (service *sdkapi) Run(ctx context.Context) error {
	service.logger.Info("running")

	err := service.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
