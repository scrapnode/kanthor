package sdkapi

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/services/command"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

func (service *sdkapi) coordinate() error {
	return service.coordinator.Receive(func(cmd string, data []byte) error {
		service.logger.Debugw("coordinating", "cmd", cmd, "data", data)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		if cmd == command.WorkspaceCredentialsCreated {
			req := &command.WorkspaceCredentialsCreatedReq{}
			if err := req.Unmarshal(data); err != nil {
				return err
			}

			for _, doc := range req.Docs {
				if err := service.authz.Grant(doc.WorkspaceId, doc.Id, RoleOwner, PermissionOwner); err != nil {
					service.logger.Errorw(err.Error(), "ws_id", doc.WorkspaceId, "wsc_id", doc.Id, "role", RoleOwner)
				}
			}

			if err := service.authz.Refresh(ctx); err != nil {
				return err
			}
		}

		if cmd == command.WorkspaceCredentialsExpired {
			req := &command.WorkspaceCredentialsExpiredReq{}
			if err := req.Unmarshal(data); err != nil {
				return err
			}

			_, err := service.uc.WorkspaceCredentials().Expire(ctx,
				&usecase.WorkspaceCredentialsExpireReq{User: req.Id, ExpiredAt: req.ExpiredAt},
			)
			if err != nil {
				service.logger.Error(err.Error())
				return err
			}
		}

		return nil
	})
}
