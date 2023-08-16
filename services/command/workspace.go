package command

import (
	"encoding/json"
	"github.com/scrapnode/kanthor/domain/entities"
)

var WorkspaceCredentialsCreated = "workspace.credentials.created"

type WorkspaceCredentialsCreatedReq struct {
	Docs []entities.WorkspaceCredentials `json:"docs"`
}

func (req *WorkspaceCredentialsCreatedReq) Marshal() ([]byte, error) {
	return json.Marshal(req)
}

func (req *WorkspaceCredentialsCreatedReq) Unmarshal(data []byte) error {
	return json.Unmarshal(data, req)
}

var WorkspaceCredentialsExpired = "workspace.credentials.expired"

type WorkspaceCredentialsExpiredReq struct {
	Id        string `json:"id"`
	ExpiredAt int64  `json:"expired_at"`
}

func (req *WorkspaceCredentialsExpiredReq) Marshal() ([]byte, error) {
	return json.Marshal(req)
}

func (req *WorkspaceCredentialsExpiredReq) Unmarshal(data []byte) error {
	return json.Unmarshal(data, req)
}
