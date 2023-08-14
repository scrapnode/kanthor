package sdkapi

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/sdkapi/docs"
	"github.com/scrapnode/kanthor/services/sdkapi/middlewares"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	validator validator.Validator,
	idempotency idempotency.Idempotency,
	coordinator coordinator.Coordinator,
	authz authorizator.Authorizator,
	uc usecase.Sdk,
) services.Service {
	logger = logger.With("service", "sdkapi")
	return &sdkapi{
		conf:        conf,
		logger:      logger,
		validator:   validator,
		coordinator: coordinator,
		idempotency: idempotency,
		authz:       authz,
		uc:          uc,
	}
}

type sdkapi struct {
	conf        *config.Config
	logger      logging.Logger
	validator   validator.Validator
	idempotency idempotency.Idempotency
	coordinator coordinator.Coordinator
	authz       authorizator.Authorizator
	uc          usecase.Sdk

	server *http.Server
}

func (service *sdkapi) Start(ctx context.Context) error {
	if err := service.authz.Connect(ctx); err != nil {
		return err
	}

	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	if err := service.idempotency.Connect(ctx); err != nil {
		return err
	}

	if err := service.coordinator.Connect(ctx); err != nil {
		return err
	}

	service.build()

	service.logger.Info("started")
	return nil
}

func (service *sdkapi) build() {
	router := gin.New()
	router.Use(gin.Recovery())
	// system routes
	router.GET("/", func(ginctx *gin.Context) {
		ginctx.JSON(http.StatusOK, gin.H{"version": service.conf.Version})
	})
	router.GET("/readiness", func(ginctx *gin.Context) {
		// @TODO: add starting up checking here
		ginctx.String(http.StatusOK, "ready")
	})
	router.GET("/liveness", func(ginctx *gin.Context) {
		ginctx.String(http.StatusOK, "live")
	})
	swagger := router.Group("/swagger")
	{
		swagger.GET("/*any", ginswagger.WrapHandler(
			swaggerfiles.Handler,
			ginswagger.PersistAuthorization(true),
			ginswagger.InstanceName(docs.SwaggerInfoSdk.InfoInstanceName),
		))
	}
	// api routes
	api := router.Group("/api")
	{
		api.Use(ginmw.UseStartup())
		api.Use(ginmw.UseIdempotency(service.logger, service.idempotency))
		api.Use(middlewares.UseAuth(service.validator, service.uc))
		api.Use(middlewares.UseAuthz(service.authz))
		api.Use(ginmw.UsePaging(service.logger, 5, 30))
		UseApplicationRoutes(
			api.Group("/application"),
			service.logger, service.validator, service.uc,
		)
		UseEndpointRoutes(
			api.Group("/application/:app_id/endpoint"),
			service.logger, service.validator, service.uc,
		)
		UseEndpointRuleRoutes(
			api.Group("/application/:app_id/endpoint/:ep_id/rule"),
			service.logger, service.validator, service.uc,
		)

		UseMessageRoutes(
			api.Group("/application/:app_id/message"),
			service.logger, service.validator, service.uc,
		)
	}

	service.server = &http.Server{
		Addr:    service.conf.SdkApi.Gateway.Httpx.Addr,
		Handler: router,
	}
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

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.authz.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *sdkapi) Run(ctx context.Context) error {
	if err := service.coordinate(); err != nil {
		return err
	}

	service.logger.Infow("running", "addr", service.conf.SdkApi.Gateway.Httpx.Addr)

	err := service.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (service *sdkapi) coordinate() error {
	// return error will deliver message again
	return service.coordinator.Receive(func(cmd *coordinator.Command) error {
		service.logger.Infow("received coordinator command", "cmd.name", cmd.Name)

		if cmd.Name == coordinator.CmdAuthzRefresh {
			if err := service.authz.Refresh(context.Background()); err != nil {
				service.logger.Error(err.Error())
				return err
			}
		}

		if cmd.Name == "workspace.credentials.create" {
			var request map[string]string
			if err := json.Unmarshal([]byte(cmd.Request), &request); err != nil {
				service.logger.Error(err.Error())
				// we don't want to process this message again because of malformed data
				return nil
			}
			if request["id"] == "" {
				service.logger.Errorw("no credentials id", "cmd", cmd.String())
				// we don't want to process this message again because of missing id
				return nil
			}
			if request["workspace_id"] == "" {
				service.logger.Errorw("no workspace id", "cmd", cmd.String())
				// we don't want to process this message again because of missing id
				return nil
			}

			err := service.authz.GrantPermissionsToRole(request["workspace_id"], RoleOwner, PermissionOwner)
			if err != nil {
				service.logger.Error(err.Error())
				return err
			}
			err = service.authz.GrantRoleToSub(request["workspace_id"], RoleOwner, request["id"])
			if err != nil {
				service.logger.Error(err.Error())
				return err
			}

			return nil
		}

		// @TODO: refresh cache

		return nil
	})
}
