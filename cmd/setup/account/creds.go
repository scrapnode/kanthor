package account

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/services/permissions"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

func creds(uc usecase.Portal, ctx context.Context, ws *entities.Workspace, out *output) error {
	ucreq := &usecase.WorkspaceCredentialsGenerateReq{
		WsId:        ws.Id,
		Name:        fmt.Sprintf("setup at %s", time.Now().UTC().Format(time.RFC3339)),
		Role:        permissions.SdkOwner,
		Permissions: permissions.SdkOwnerPermissions,
	}
	if err := ucreq.Validate(); err != nil {
		return err
	}
	ucres, err := uc.WorkspaceCredentials().Generate(ctx, ucreq)
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
	t.AppendHeader(table.Row{"key", "secret", "encoded"})

	encoded := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pass)))
	t.AppendRow([]interface{}{user, pass, encoded})
	return t.Render()
}
