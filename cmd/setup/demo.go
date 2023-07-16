package setup

import (
	"context"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/controlplane/permissions"
	"github.com/scrapnode/kanthor/services/ioc"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	"os"
	"strings"
	"time"
)

func Demo(conf *config.Config, logger logging.Logger, owner string, verbose bool) error {
	uc, err := ioc.InitializeControlplaneUsecase(conf, logger)
	if err != nil {
		return err
	}

	authz := authorizator.New(&conf.Controlplane.Authorizator, logger)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	if err := uc.Connect(ctx); err != nil {
		return err
	}
	if err := authz.Connect(ctx); err != nil {
		return err
	}
	defer func() {
		if err := uc.Disconnect(ctx); err != nil {
			logger.Error(err)
		}

		if err := authz.Disconnect(ctx); err != nil {
			logger.Error(err)
		}
	}()

	acc := &authenticator.Account{Sub: owner}

	// @TODO: add transaction from here
	project, err := uc.Project().SetupDefault(ctx, &usecase.ProjectSetupDefaultReq{Account: acc})
	if err != nil {
		return err
	}

	policies := permissions.PoliciesOfRoleInWorkspace(permissions.RoleOwner, project.WorkspaceId)
	if err := authz.AddPolicies(policies); err != nil {
		return err
	}
	if err := authz.Grant(acc.Sub, permissions.RoleOwner, project.WorkspaceId); err != nil {
		return err
	}

	data, err := uc.Project().SetupDemo(ctx, &usecase.ProjectSetupDemoReq{
		Account:     acc,
		WorkspaceId: project.WorkspaceId,
	})
	if err != nil {
		return err
	}

	if verbose {
		t := table.NewWriter()
		t.AppendHeader(table.Row{"name", "value"})
		t.AppendRow([]interface{}{"workspace_id", project.WorkspaceId})
		t.AppendRow([]interface{}{"workspace_tier", project.WorkspaceTier})
		t.AppendRow([]interface{}{"application_id", data.ApplicationId})
		t.AppendRow([]interface{}{"endpoint_ids", strings.Join(data.EndpointIds, ", ")})
		t.AppendRow([]interface{}{"endpoint_rules_ids", strings.Join(data.EndpointRuleIds, ", ")})
		t.SetOutputMirror(os.Stdout)
		t.Render()
	}

	return nil
}
