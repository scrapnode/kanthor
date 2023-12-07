package rest

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

func RegisterWorkspaceResolver(uc usecase.Sdk) func(ctx context.Context, id string) (*entities.Workspace, error) {
	return func(ctx context.Context, id string) (*entities.Workspace, error) {
		in := &usecase.WorkspaceGetIn{Id: id}
		if err := in.Validate(); err != nil {
			return nil, err
		}

		out, err := uc.Workspace().Get(ctx, in)
		if err != nil {
			return nil, err
		}
		return out.Workspace, nil
	}
}

func RegisterBasicAuthenticate(uc usecase.Sdk) authenticator.Authenticate {
	return func(ctx context.Context, request *authenticator.Request) (*authenticator.Account, error) {
		bytes, err := base64.StdEncoding.DecodeString(request.Credentials)
		if err != nil {
			return nil, err
		}
		userpass := strings.Split(string(bytes), ":")
		if len(userpass) != 2 {
			return nil, fmt.Errorf("SDK.AUTHORIZATION.MALFORMED_CREDENTIALS")
		}

		in := &usecase.WorkspaceAuthenticateIn{User: userpass[0], Pass: userpass[1]}
		if err := in.Validate(); err != nil {
			return nil, err
		}

		out, err := uc.Workspace().Authenticate(ctx, in)
		if err != nil {
			return nil, err
		}

		acc := &authenticator.Account{
			Sub:  out.Credentials.Id,
			Name: out.Credentials.Name,
		}
		return acc, nil
	}
}
