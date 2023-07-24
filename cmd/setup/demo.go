package setup

import (
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/scrapnode/kanthor/config"
	demodata "github.com/scrapnode/kanthor/data/demo"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/pipeline"
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
		t.SetOutputMirror(os.Stdout)
		style := table.StyleDefault
		style.Format.Header = text.FormatDefault
		t.SetStyle(style)
		t.AppendHeader(table.Row{"WS - TIER", fmt.Sprintf("%s - %s", cpdata.WorkspaceId, cpdata.WorkspaceTier)})

		for _, appId := range cpdata.ApplicationIds {
			t.AppendRow([]interface{}{"app_id", appId})

			sub := "-"
			token := fmt.Sprintf("unable to generate access token for [%s]", appId)
			if authdata, ok := dpdata[appId]; ok {
				sub = authdata.Sub
				token = authdata.Token
			}
			t.AppendRow([]interface{}{"app_sub", sub})
			t.AppendRow([]interface{}{"app_sub_token", token})
			t.AppendSeparator()
		}
		t.Render()
	}

	return nil
}

func demoControlplane(
	conf *config.Config,
	logger logging.Logger,
	owner, input string,
) (*controlplaneuc.ProjectSetupDemoRes, error) {
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
	ctx = context.WithValue(ctx, authenticator.CtxAuthAccount, acc)

	project, err := demoControlplaneProjectDefault(acc, uc, ctx)
	if err != nil {
		return nil, err
	}

	data, err := demoControlplaneProjectData(acc, uc, ctx, project, bytes)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func demoControlplaneProjectDefault(
	acc *authenticator.Account,
	uc controlplaneuc.Controlplane,
	ctx context.Context,
) (*controlplaneuc.ProjectSetupDefaultRes, error) {
	pipe := pipeline.Chain(pipeline.UseValidation())
	run := pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = uc.Project().SetupDefault(ctx, request.(*controlplaneuc.ProjectSetupDefaultReq))
		return
	})

	request := &controlplaneuc.ProjectSetupDefaultReq{
		Account:       acc,
		WorkspaceName: constants.DefaultWorkspaceName,
		WorkspaceTier: constants.DefaultWorkspaceTier,
	}
	response, err := run(ctx, request)

	if err != nil {
		return nil, err
	}
	return response.(*controlplaneuc.ProjectSetupDefaultRes), nil
}

func demoControlplaneProjectData(
	acc *authenticator.Account,
	uc controlplaneuc.Controlplane,
	ctx context.Context,
	project *controlplaneuc.ProjectSetupDefaultRes,
	bytes []byte,
) (*controlplaneuc.ProjectSetupDemoRes, error) {
	pipe := pipeline.Chain(pipeline.UseValidation())
	run := pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = uc.Project().SetupDemo(ctx, request.(*controlplaneuc.ProjectSetupDemoReq))
		return
	})

	entities, err := demodata.Project(project.WorkspaceId, acc, bytes)
	if err != nil {
		return nil, err
	}

	request := &controlplaneuc.ProjectSetupDemoReq{
		Account:       acc,
		WorkspaceId:   project.WorkspaceId,
		Applications:  entities.Applications,
		Endpoints:     entities.Endpoints,
		EndpointRules: entities.EndpointRules,
	}
	response, err := run(ctx, request)

	if err != nil {
		return nil, err
	}

	return response.(*controlplaneuc.ProjectSetupDemoRes), nil
}

func demoDataplane(conf *config.Config, logger logging.Logger, appIds []string) (map[string]*dataplaneuc.AppCredsCreateRes, error) {
	maps := map[string]*dataplaneuc.AppCredsCreateRes{}

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
		req := &dataplaneuc.AppCredsCreateReq{
			AppId:       appId,
			Role:        permissions.Admin,
			Permissions: permissions.AdminPermission,
		}
		if res, err := uc.AppCreds().Create(ctx, req); err == nil {
			maps[appId] = res
		}
	}

	return maps, nil
}
