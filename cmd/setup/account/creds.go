package account

import (
	"context"
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	"github.com/scrapnode/kanthor/services/command"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
)

func creds(coord coordinator.Coordinator, uc usecase.Portal, ctx context.Context, ws *entities.Workspace, out *output) error {
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

	out.AddStdout(credsOutput(ucres.Credentials.Id, ucres.Password))
	out.AddJson("credentials", credentials{Username: ucres.Credentials.Id, Password: ucres.Password})

	return nil
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func credsOutput(user, pass string) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"key", "secret"})
	t.AppendRow([]interface{}{user, pass})
	return t.Render()
}
