package account

import (
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/services/command"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
	"os"
	"time"
)

func creds(coord coordinator.Coordinator, uc usecase.Portal, ctx context.Context, ws *entities.Workspace, withCreds bool) error {
	if !withCreds {
		return nil
	}

	ucreq := &usecase.WorkspaceCredentialsGenerateReq{
		WorkspaceId: ws.Id,
		Name:        fmt.Sprintf("setup at %s", time.Now().UTC().Format(time.RFC3339)),
	}
	if err := validator.New().Struct(ucreq); err != nil {
		return err
	}
	ucres, err := uc.WorkspaceCredentials().Generate(ctx, ucreq)
	if err != nil {
		return err
	}

	err = coord.Send(
		ctx,
		command.WorkspaceCredentialsCreated,
		&command.WorkspaceCredentialsCreatedReq{Docs: []entities.WorkspaceCredentials{*ucres.Credentials}},
	)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"ws", "key", "secret"})
	t.AppendRow([]interface{}{ucres.Credentials.WorkspaceId, ucres.Credentials.Id, ucres.Password})
	t.SetOutputMirror(os.Stdout)
	t.Render()
	return nil
}
