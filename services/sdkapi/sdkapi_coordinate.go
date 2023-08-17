package sdkapi

import (
	"context"
	"github.com/scrapnode/kanthor/services/command"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"time"
)

func (service *sdkapi) coordinate() error {
	return service.coordinator.Receive(func(cmd string, data []byte) error {
		service.logger.Debugw("coordinating", "cmd", cmd, "data", data)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		if cmd == command.WorkspaceCredentialsCreated {
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
