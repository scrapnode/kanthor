package entities

import "github.com/scrapnode/kanthor/project"

var (
	TableWs  = project.NameWithoutTier("workspace")
	TableWsc = project.NameWithoutTier("workspace_credentials")
	TableApp = project.NameWithoutTier("application")
	TableEp  = project.NameWithoutTier("endpoint")
	TableEpr = project.NameWithoutTier("endpoint_rule")
	TableMsg = project.NameWithoutTier("message")
	TableReq = project.NameWithoutTier("request")
	TableRes = project.NameWithoutTier("response")
	TableAtt = project.NameWithoutTier("attempt")
)
