package sdkapi

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	RoleOwner = "sdk.owner"
)

var PermissionOwner = []authorizator.Permission{
	{Object: "/api/account/me", Action: "GET"},

	{Object: "/api/application", Action: "POST"},
	{Object: "/api/application/:app_id", Action: "PUT"},
	{Object: "/api/application/:app_id", Action: "DELETE"},
	{Object: "/api/application", Action: "GET"},
	{Object: "/api/application/:app_id", Action: "GET"},

	{Object: "/api/application/:app_id/endpoint", Action: "POST"},
	{Object: "/api/application/:app_id/endpoint/:ep_id", Action: "PUT"},
	{Object: "/api/application/:app_id/endpoint/:ep_id", Action: "DELETE"},
	{Object: "/api/application/:app_id/endpoint", Action: "GET"},
	{Object: "/api/application/:app_id/endpoint/:ep_id", Action: "GET"},

	{Object: "/api/application/:app_id/endpoint/:ep_id/rule", Action: "POST"},
	{Object: "/api/application/:app_id/endpoint/:ep_id/rule/:epr_id", Action: "PUT"},
	{Object: "/api/application/:app_id/endpoint/:ep_id/rule/:epr_id", Action: "DELETE"},
	{Object: "/api/application/:app_id/endpoint/:ep_id/rule", Action: "GET"},
	{Object: "/api/application/:app_id/endpoint/:ep_id/rule/:epr_id", Action: "GET"},

	{Object: "/api/application/:app_id/message", Action: "PUT"},
}
