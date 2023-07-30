package do

import (
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/data/interchange"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/ioc"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func NewImport(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:  "import",
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

			// prepare the data
			bytes, err := os.ReadFile(input)
			if err != nil {
				return err
			}
			in, err := interchange.Demo(ownerId, bytes)
			if err != nil {
				return err
			}

			req, res, err := importDo(conf, logger, in)
			if err != nil {
				return err
			}

			importReport(in, req, res, verbose)
			return nil
		},
	}

	command.Flags().StringArrayP("auto-generate", "", []string{}, "--auto-generate=workspace_credentials | auto generate some value that could not be exported & imported")

	return command
}

func importDo(conf *config.Config, logger logging.Logger, in *interchange.Interchange) (*usecase.WorkspaceImportReq, *usecase.WorkspaceImportRes, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*2)
	defer cancel()

	uc, err := ioc.InitializePortalUsecase(conf, logger)
	if err != nil {
		return nil, nil, err
	}
	if err := uc.Connect(ctx); err != nil {
		return nil, nil, err
	}
	defer func() {
		if err := uc.Disconnect(ctx); err != nil {
			logger.Error(err)
		}
	}()

	req := &usecase.WorkspaceImportReq{}
	for _, workspace := range in.Workspaces {
		req.Workspaces = append(req.Workspaces, *workspace.Workspace)

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
		return nil, nil, err
	}

	return req, res, nil
}

func importReport(in *interchange.Interchange, req *usecase.WorkspaceImportReq, res *usecase.WorkspaceImportRes, verbose bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	style := table.StyleDefault
	style.Format.Header = text.FormatDefault
	t.SetStyle(style)
	title := fmt.Sprintf(
		"Summary: %d/%d worksapces, %d tiers, %d/%d apps, %d/%d endpoints, %d/%d rules",
		len(res.WorkspaceIds), len(req.Workspaces),
		len(res.WorkspaceTierIds),
		len(res.ApplicationIds), len(req.Applications),
		len(res.EndpointIds), len(req.Endpoints),
		len(res.EndpointRuleIds), len(req.EndpointRules),
	)
	t.SetTitle(title)

	for _, ws := range in.Workspaces {
		t.AppendRow([]interface{}{"ws_id", ws.Id})
		t.AppendSeparator()

		if verbose {
			for _, app := range ws.Applications {
				t.AppendRow([]interface{}{"app_id", app.Id})
				for _, ep := range app.Endpoints {
					t.AppendRow([]interface{}{"ep_id", ep.Id})
					for _, epr := range ep.Rules {
						t.AppendRow([]interface{}{"epr_id", epr.Id})
					}
				}
				t.AppendSeparator()
			}
		}

		t.AppendSeparator()
	}

	t.Render()
}
