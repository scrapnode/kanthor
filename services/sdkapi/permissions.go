package sdkapi

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	RoleOwner = "sdk.owner"
)

var PermissionOwner = []authorizator.Permission{
	{"/api/application", "POST"},
	{"/api/application/:app_id", "PUT"},
	{"/api/application/:app_id", "DELETE"},
	{"/api/application", "GET"},
	{"/api/application/:app_id", "GET"},

	{"/api/application/:app_id/endpoint", "POST"},
	{"/api/application/:app_id/endpoint/:ep_id", "PUT"},
	{"/api/application/:app_id/endpoint/:ep_id", "DELETE"},
	{"/api/application/:app_id/endpoint", "GET"},
	{"/api/application/:app_id/endpoint/:ep_id", "GET"},

	{"/api/application/:app_id/endpoint/:ep_id/rule", "POST"},
	{"/api/application/:app_id/endpoint/:ep_id/rule/:epr_id", "PUT"},
	{"/api/application/:app_id/endpoint/:ep_id/rule/:epr_id", "DELETE"},
	{"/api/application/:app_id/endpoint/:ep_id/rule", "GET"},
	{"/api/application/:app_id/endpoint/:ep_id/rule/:epr_id", "GET"},

	{"/api/application/:app_id/message", "PUT"},
}
