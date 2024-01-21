package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type ApplicationListMessageRes struct {
	Data []Message `json:"data"`
} // @name ApplicationListMessageRes

// UseApplicationListMessage
// @Tags		application
// @Router		/application/{app_id}/message		[get]
// @Param		_limit								query		int				false	"limit returning records" 	default(10)
// @Param		_start								query		int64			false	"starting time to scan in milliseconds" example(1669914060000)
// @Param		_end								query		int64			false	"ending time to scan in milliseconds" 	example(1985533260000)
// @Success		200									{object}	ApplicationListMessageRes
// @Failure		default								{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseApplicationListMessage(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req gateway.Query
		if err := ginctx.BindQuery(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)

		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)
		in := &usecase.ApplicationListMessageIn{
			ScanningQuery: entities.ScanningQueryFromGatewayQuery(&req, service.infra.Timer),
			WsId:          ws.Id,
			AppId:         ginctx.Param("app_id"),
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Application().ListMessage(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &ApplicationListMessageRes{Data: make([]Message, 0)}
		for _, doc := range out.Data {
			res.Data = append(res.Data, *ToMessage(&doc))
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
