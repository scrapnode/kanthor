package do

import (
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/data/interchange"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/services/command"
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
			sub, err := cmd.Flags().GetString("account-sub")
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString("account-name")
			if err != nil {
				return err
			}
			if name == "" {
				name = sub
			}
			needs, err := cmd.Flags().GetStringArray("auto-generate")
			if err != nil {
				return err
			}
			input := args[0]

			ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*2)
			defer cancel()

			meter, err := metrics.NewNoop(nil, logger)
			if err != nil {
				return err
			}
			uc, err := ioc.InitializePortalUsecase(conf, logger, meter)
			if err != nil {
				return err
			}
			if err := uc.Connect(ctx); err != nil {
				return err
			}

			coord, err := coordinator.New(&conf.Coordinator, logger)
			if err != nil {
				return err
			}
			if err := coord.Connect(ctx); err != nil {
				return err
			}

			defer func() {
				if err := uc.Disconnect(ctx); err != nil {
					logger.Error(err)
				}

				if err := coord.Disconnect(ctx); err != nil {
					logger.Error(err)
				}

			}()

			// prepare the data
			acc := &authenticator.Account{Sub: sub, Name: name}
			ctx = context.WithValue(ctx, authenticator.CtxAcc, acc)
			bytes, err := os.ReadFile(input)
			if err != nil {
				return err
			}
			in, err := interchange.Demo(acc, bytes)
			if err != nil {
				return err
			}
			req := importPrepareRequest(in)

			res, err := uc.Workspace().Import(ctx, req)
			if err != nil {
				return err
			}

			metadata := importAutoGenerate(uc, coord, res, ctx, needs)

			importReport(in, req, res, metadata, verbose)
			return nil
		},
	}

	command.Flags().StringArrayP("auto-generate", "", []string{}, "--auto-generate=workspace_credentials | auto generate some value that could not be exported & imported")

	return command
}

func importPrepareRequest(in *interchange.Interchange) *usecase.WorkspaceImportReq {
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

	return req
}

func importAutoGenerate(
	uc usecase.Portal,
	coord coordinator.Coordinator,
	res *usecase.WorkspaceImportRes,
	ctx context.Context,
	needs []string,
) entities.Metadata {
	meta := entities.Metadata{}

	for _, need := range needs {
		if need == "workspace_credentials" {
			importAutoGenerateWorkspaceCredentials(uc, coord, res, ctx, meta)
		}
	}

	return meta
}

func importAutoGenerateWorkspaceCredentials(
	uc usecase.Portal,
	coord coordinator.Coordinator,
	res *usecase.WorkspaceImportRes,
	ctx context.Context,
	meta entities.Metadata,
) {
	for _, wsId := range res.WorkspaceIds {
		cred, err := uc.WorkspaceCredentials().Generate(
			ctx,
			&usecase.WorkspaceCredentialsGenerateReq{WorkspaceId: wsId, Count: 1},
		)

		if err != nil {
			meta.Set(wsId, fmt.Sprintf("error: %s", err.Error()))
			continue
		}

		err = coord.Send(
			ctx,
			command.WorkspaceCredentialsCreated,
			&command.WorkspaceCredentialsCreatedReq{Docs: cred.Credentials},
		)
		if err != nil {
			meta.Set(wsId, fmt.Sprintf("error: %s", err.Error()))
			continue
		}

		token := fmt.Sprintf("%s:%s", cred.Credentials[0].Id, cred.Passwords[cred.Credentials[0].Id])
		meta.Set(wsId, fmt.Sprintf("user:pass [%s]", token))
	}
}

func importReport(
	in *interchange.Interchange,
	req *usecase.WorkspaceImportReq,
	res *usecase.WorkspaceImportRes,
	metadata entities.Metadata,
	verbose bool,
) {
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

	t.AppendHeader([]interface{}{"key", "value", "meta"})
	for _, ws := range in.Workspaces {
		t.AppendRow([]interface{}{"ws_id", ws.Id, metadata.Get(ws.Id)})
		t.AppendSeparator()

		if verbose {
			for _, app := range ws.Applications {
				t.AppendRow([]interface{}{"app_id", app.Id, metadata.Get(app.Id)})
				for _, ep := range app.Endpoints {
					t.AppendRow([]interface{}{"ep_id", ep.Id, metadata.Get(ep.Id)})
					for _, epr := range ep.Rules {
						t.AppendRow([]interface{}{"epr_id", epr.Id, metadata.Get(epr.Id)})
					}
				}
				t.AppendSeparator()
			}
		}

		t.AppendSeparator()
	}

	t.Render()
}
