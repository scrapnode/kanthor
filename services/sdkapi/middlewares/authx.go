package middlewares

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/pkg/sender"
	sdkuc "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

var (
	HeaderAuthScheme  = "kanthor-auth-scheme"
	AuthSchemeForward = "forward"
)

type PortalAuth struct {
	Error         string                  `json:"error"`
	Account       *authenticator.Account  `json:"account"`
	Workspace     *entities.Workspace     `json:"workspace"`
	WorkspaceTier *entities.WorkspaceTier `json:"workspace_tier"`
}

func UseAuthx(
	conf config.SdkApi,
	logger logging.Logger,
	validator validator.Validator,
	authz authorizator.Authorizator,
	uc sdkuc.Sdk,
) gin.HandlerFunc {
	send := sender.Rest(sender.DefaultConfig, logger)

	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)

		// portal authentication & authorization
		if ginctx.Request.Header.Get(HeaderAuthScheme) == AuthSchemeForward {
			if conf.PortalConnection.Account == "" {
				ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no portal connection endpoint was configured"})
				return
			}

			req := &sender.Request{
				Method:  http.MethodGet,
				Headers: ginctx.Request.Header.Clone(),
				Uri:     conf.PortalConnection.Account,
			}
			res, err := send(req)
			if err != nil {
				logger.Error(err)
				ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "remote auth is failed"})
				return
			}

			var auth PortalAuth
			if err := json.Unmarshal(res.Body, &auth); err != nil {
				logger.Error(err)
				ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unable to parse auth response from remote server"})
				return
			}

			if auth.Error != "" {
				ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": auth.Error})
				return
			}

			if auth.Account == nil {
				ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no account"})
				return
			}
			ctx = context.WithValue(ctx, authenticator.CtxAcc, auth.Account)
			if auth.Workspace == nil {
				ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no workspace"})
				return
			}
			ctx = context.WithValue(ctx, authorizator.CtxWs, auth.Workspace)
			if auth.WorkspaceTier == nil {
				ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no workspace tier"})
				return
			}
			ctx = context.WithValue(ctx, authorizator.CtxWst, auth.WorkspaceTier)

			ginctx.Set(gateway.KeyCtx, ctx)
			ginctx.Next()
			return
		}

		// sdk authentication & authorization
		credentials := ginctx.Request.Header.Get(authenticator.HeaderAuth)
		ctx, err := sdkauth(validator, uc, ctx, credentials)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unable to verify your credentials"})
			return
		}
		acc := ctx.Value(authenticator.CtxAcc).(*authenticator.Account)
		ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
		obj := ginctx.FullPath() // form of /application/:app_id
		act := ginctx.Request.Method
		ok, err := authz.Enforce(acc.Sub, ws.Id, obj, act)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unable to enforce permission"})
			return
		}
		if !ok {
			ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you have no permission to perform the action"})
			return
		}

		ginctx.Set(gateway.KeyCtx, ctx)
		ginctx.Next()
	}
}

func sdkauth(
	validator validator.Validator,
	uc sdkuc.Sdk,
	ctx context.Context,
	credentials string,
) (context.Context, error) {
	user, pass, err := authenticator.ParseBasicCredentials(credentials)
	if err != nil {
		return ctx, err
	}

	req := &sdkuc.WorkspaceCredentialsAuthenticateReq{User: user, Pass: pass}
	if err := validator.Struct(req); err != nil {
		return ctx, err
	}

	res, err := uc.WorkspaceCredentials().Authenticate(ctx, req)
	if err != nil {
		return ctx, err
	}

	acc := &authenticator.Account{Sub: user, Name: res.WorkspaceCredentials.Name}
	ctx = context.WithValue(ctx, authenticator.CtxAcc, acc)
	ctx = context.WithValue(ctx, authorizator.CtxWs, res.Workspace)
	ctx = context.WithValue(ctx, authorizator.CtxWst, res.WorkspaceTier)
	return ctx, nil
}
