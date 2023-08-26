package portalapi

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	RoleOwner = "portal.owner"
)

var PermissionOwner = []authorizator.Permission{
	{Object: "/api/workspace/me", Action: "GET"},
	{Object: "/api/workspace/me", Action: "PUT"},

	{Object: "/api/workspace/me/credentials", Action: "POST"},
	{Object: "/api/workspace/me/credentials", Action: "GET"},
	{Object: "/api/workspace/me/credentials/:wsc_id", Action: "PUT"},
	{Object: "/api/workspace/me/credentials/:wsc_id", Action: "GET"},
	{Object: "/api/workspace/me/credentials/:wsc_id/expiration", Action: "PUT"},
}
