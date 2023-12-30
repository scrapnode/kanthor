package permissions

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	PortalOwner = "portal.owner"
)

var PortalOwnerPermissions = []authorizator.Permission{
	{Object: "/api/workspace", Action: "GET"},
	{Object: "/api/workspace", Action: "POST"},
	{Object: "/api/workspace/:ws_id", Action: "GET"},
	{Object: "/api/workspace/:ws_id", Action: "PATCH"},

	{Object: "/api/credentials", Action: "POST"},
	{Object: "/api/credentials", Action: "GET"},
	{Object: "/api/credentials/:wsc_id", Action: "PATCH"},
	{Object: "/api/credentials/:wsc_id", Action: "GET"},
	{Object: "/api/credentials/:wsc_id/expiration", Action: "PUT"},
	{Object: "/api/application/:app_id/message", Action: "GET"},
	{Object: "/api/application/:app_id/message/:msg_id", Action: "GET"},
	{Object: "/api/endpoint/:ep_id/message", Action: "GET"},
	{Object: "/api/endpoint/:ep_id/message/:msg_id", Action: "GET"},
}
