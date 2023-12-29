package permissions

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	PortalOwner = "portal.owner"
)

var PortalOwnerPermissions = []authorizator.Permission{
	{Object: "/api/workspace/me", Action: "GET"},
	{Object: "/api/workspace/me", Action: "PATCH"},

	{Object: "/api/credentials", Action: "POST"},
	{Object: "/api/credentials", Action: "GET"},
	{Object: "/api/credentials/:wsc_id", Action: "PATCH"},
	{Object: "/api/credentials/:wsc_id", Action: "GET"},
	{Object: "/api/credentials/:wsc_id/expiration", Action: "PUT"},
	{Object: "/api/message", Action: "GET"},
	{Object: "/api/message/:msg_id", Action: "GET"},
	{Object: "/api/request", Action: "GET"},
	{Object: "/api/request/:req_id", Action: "GET"},
	{Object: "/api/response", Action: "GET"},
	{Object: "/api/response/:res_id", Action: "GET"},
}
