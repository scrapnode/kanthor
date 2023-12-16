package permissions

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	PortalOwner = "portal.owner"
)

var PortalOwnerPermissions = []authorizator.Permission{
	{Object: "/api/workspace/me", Action: "GET"},
	{Object: "/api/workspace/me", Action: "PUT"},

	{Object: "/api/workspace/me/credentials", Action: "POST"},
	{Object: "/api/workspace/me/credentials", Action: "GET"},
	{Object: "/api/workspace/me/credentials/:wsc_id", Action: "PATCH"},
	{Object: "/api/workspace/me/credentials/:wsc_id", Action: "GET"},
	{Object: "/api/workspace/me/credentials/:wsc_id/expiration", Action: "PUT"},
}
