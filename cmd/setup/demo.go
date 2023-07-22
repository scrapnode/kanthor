package setup

import (
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/scrapnode/kanthor/config"
	demodata "github.com/scrapnode/kanthor/data/demo"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/dataplane/permissions"
	"github.com/scrapnode/kanthor/services/ioc"
	controlplaneuc "github.com/scrapnode/kanthor/usecases/controlplane"
	dataplaneuc "github.com/scrapnode/kanthor/usecases/dataplane"
	"os"
	"time"
)

func demo(conf *config.Config, logger logging.Logger, owner, input string, verbose bool) error {
	cpdata, err := demoControlplane(conf, logger, owner, input)
	if err != nil {
		return err
	}

	dpdata, err := demoDataplane(conf, logger, cpdata.ApplicationIds)
	if err != nil {
		return err
	}

	if verbose {
		t := table.NewWriter()
		t.AppendHeader(table.Row{"workspace_id", "workspace_tier", "application_id", "auth_sub", "auth_token"})
		for _, appId := range cpdata.ApplicationIds {
			sub := "-"
			token := fmt.Sprintf("unable to generate access token for [%s]", appId)
			if authdata, ok := dpdata[appId]; ok {
				sub = authdata.Sub
				token = authdata.Token
			}
			t.AppendRow([]interface{}{cpdata.WorkspaceTier, cpdata.WorkspaceTier, appId, sub, token})
		}
		t.SetOutputMirror(os.Stdout)
		t.Render()
	}

	return nil
}

func demoControlplane(conf *config.Config, logger logging.Logger, owner, input string) (*controlplaneuc.ProjectSetupDemoRes, error) {
	bytes, err := os.ReadFile(input)
	if err != nil {
		return nil, err
	}

	uc, err := ioc.InitializeControlplaneUsecase(conf, logger)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	if err := uc.Connect(ctx); err != nil {
		return nil, err
	}
	defer func() {
		if err := uc.Disconnect(ctx); err != nil {
			logger.Error(err)
		}
	}()

	acc := &authenticator.Account{Sub: owner}
	project, err := uc.Project().SetupDefault(ctx, &controlplaneuc.ProjectSetupDefaultReq{Account: acc})
	if err != nil {
		return nil, err
	}

	entities, err := demodata.Project(project.WorkspaceId, bytes)
	if err != nil {
		return nil, err
	}
	return uc.Project().SetupDemo(ctx, &controlplaneuc.ProjectSetupDemoReq{
		Account:       acc,
		WorkspaceId:   project.WorkspaceId,
		Applications:  entities.Applications,
		Endpoints:     entities.Endpoints,
		EndpointRules: entities.EndpointRules,
	})
}

func demoDataplane(conf *config.Config, logger logging.Logger, appIds []string) (map[string]*dataplaneuc.ApplicationGenTokenRes, error) {
	maps := map[string]*dataplaneuc.ApplicationGenTokenRes{}

	uc, err := ioc.InitializeDataplaneUsecase(conf, logger)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	if err := uc.Connect(ctx); err != nil {
		return nil, err
	}
	defer func() {
		if err := uc.Disconnect(ctx); err != nil {
			logger.Error(err)
		}
	}()

	for _, appId := range appIds {
		req := &dataplaneuc.ApplicationGenTokenReq{
			Id:          appId,
			Role:        permissions.Admin,
			Permissions: permissions.AdminPermission,
		}
		if res, err := uc.Application().GenToken(ctx, req); err == nil {
			maps[appId] = res
		}
	}

	return maps, nil
}
