package interchange

import (
	"encoding/json"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"time"
)

func Demo(cryptor cryptography.Cryptography, ownerId string, bytes []byte) (*Interchange, error) {
	var in Interchange
	if err := json.Unmarshal(bytes, &in); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	for i, workspace := range in.Workspaces {
		workspace.GenId()
		workspace.SetAT(now)
		workspace.OwnerId = ownerId
		workspace.ModifiedBy = ownerId

		tier := workspace.Tier
		tier.GenId()
		tier.SetAT(now)
		tier.WorkspaceId = workspace.Id
		tier.ModifiedBy = ownerId
		workspace.Tier = tier

		credentials := entities.WorkspaceCredentials{}
		credentials.GenId()
		credentials.SetAT(now)
		credentials.WorkspaceId = workspace.Id
		credentials.ModifiedBy = ownerId
		hashed, err := cryptor.KDF().StringHash(credentials.Id)
		if err != nil {
			return nil, err
		}
		credentials.Hash = hashed

		workspace.Credentials = append(workspace.Credentials, credentials)

		for j, application := range workspace.Applications {
			application.GenId()
			application.SetAT(now)
			application.ModifiedBy = ownerId

			for k, endpoint := range application.Endpoints {
				endpoint.GenId()
				endpoint.SetAT(now)
				endpoint.AppId = application.Id
				endpoint.ModifiedBy = ownerId

				for h, rule := range endpoint.Rules {
					rule.GenId()
					rule.SetAT(now)
					rule.EndpointId = endpoint.Id
					rule.ModifiedBy = ownerId

					endpoint.Rules[h] = rule
				}

				application.Endpoints[k] = endpoint
			}

			application.WorkspaceId = workspace.Id
			workspace.Applications[j] = application
		}

		in.Workspaces[i] = workspace
	}

	return &in, nil
}
