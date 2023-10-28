package entities

import "github.com/scrapnode/kanthor/namespace"

var (
	TableWs  = namespace.NameWithoutTier("workspace")
	TableWsc = namespace.NameWithoutTier("workspace_credentials")
	TableApp = namespace.NameWithoutTier("application")
	TableEp  = namespace.NameWithoutTier("endpoint")
	TableEpr = namespace.NameWithoutTier("endpoint_rule")
	TableMsg = namespace.NameWithoutTier("message")
	TableReq = namespace.NameWithoutTier("request")
	TableRes = namespace.NameWithoutTier("response")
	TableAtt = namespace.NameWithoutTier("attempt")
)
