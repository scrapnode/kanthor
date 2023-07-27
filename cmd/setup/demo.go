package setup

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/data/interchange"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/ioc"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func NewDemo(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:  "demo",
		Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}
			ownerId, err := cmd.Flags().GetString("account-sub")
			if err != nil {
				return err
			}
			input := args[0]

			return demo(conf, logger, input, ownerId, verbose)
		},
	}

	return command
}

func demo(conf *config.Config, logger logging.Logger, input, ownerId string, verbose bool) error {
	bytes, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	cryptor, err := cryptography.New(&conf.Cryptography)
	if err != nil {
		return err
	}

	in, err := interchange.Demo(cryptor, ownerId, bytes)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*2)
	defer cancel()

	uc, err := ioc.InitializePortalUsecase(conf, logger)
	if err != nil {
		return err
	}
	if err := uc.Connect(ctx); err != nil {
		return err
	}
	defer func() {
		if err := uc.Disconnect(ctx); err != nil {
			logger.Error(err)
		}
	}()

	req := &usecase.WorkspaceImportReq{}
	for _, workspace := range in.Workspaces {
		req.Workspaces = append(req.Workspaces, *workspace.Workspace)
		req.WorkspaceTiers = append(req.WorkspaceTiers, *workspace.Tier)
		req.WorkspaceCredentials = append(req.WorkspaceCredentials, workspace.Credentials...)

		for _, app := range workspace.Applications {
			req.Applications = append(req.Applications, *app.Application)

			for _, ep := range app.Endpoints {
				req.Endpoints = append(req.Endpoints, *ep.Endpoint)

				for _, epr := range ep.Rules {
					req.EndpointRules = append(req.EndpointRules, *epr.EndpointRule)
				}
			}
		}
	}

	res, err := uc.Workspace().Import(ctx, req)
	if err != nil {
		return err
	}

	if verbose {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		style := table.StyleDefault
		style.Format.Header = text.FormatDefault
		t.SetStyle(style)
		count := len(res.WorkspaceIds) + len(res.WorkspaceTierIds) + len(res.WorkspaceCredentialsIds) + len(res.ApplicationIds) + len(res.EndpointIds) + len(res.EndpointRuleIds)
		t.SetTitle(fmt.Sprintf("Import Count: %d items", count))

		for _, workspace := range in.Workspaces {
			credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", workspace.Credentials[0].Id, workspace.Credentials[0].Id)))
			t.AppendHeader(table.Row{"WS - TIER - Credentials", fmt.Sprintf("%s - %s - %s", workspace.Id, workspace.Tier.Name, credentials)})

			for _, app := range workspace.Applications {
				t.AppendRow([]interface{}{"app_id", app.Id})
			}

			t.AppendSeparator()
		}

		t.Render()
	}

	return nil
}
