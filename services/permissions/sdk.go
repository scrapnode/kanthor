package permissions

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	SdkOwner = "sdk.owner"
)

var SdkOwnerPermissions = []authorizator.Permission{
	{Object: "/api/account/me", Action: "GET"},

	{Object: "/api/application", Action: "POST"},
	{Object: "/api/application/:app_id", Action: "PATCH"},
	{Object: "/api/application/:app_id", Action: "DELETE"},
	{Object: "/api/application", Action: "GET"},
	{Object: "/api/application/:app_id", Action: "GET"},

	{Object: "/api/endpoint", Action: "POST"},
	{Object: "/api/endpoint/:ep_id", Action: "PATCH"},
	{Object: "/api/endpoint/:ep_id", Action: "DELETE"},
	{Object: "/api/endpoint", Action: "GET"},
	{Object: "/api/endpoint/:ep_id", Action: "GET"},
	{Object: "/api/endpoint/:ep_id/secret", Action: "GET"},

	{Object: "/api/rule", Action: "POST"},
	{Object: "/api/rule/:epr_id", Action: "PATCH"},
	{Object: "/api/rule/:epr_id", Action: "DELETE"},
	{Object: "/api/rule", Action: "GET"},
	{Object: "/api/rule/:epr_id", Action: "GET"},

	{Object: "/api/message", Action: "POST"},
}
