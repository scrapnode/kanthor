package portalapi

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	RoleOwner = "owner"
)

var PermissionOwner = append([]authorizator.Permission{
	{"/api/workspace/me", "GET"},

	{"/api/workspace/me/credentials", "POST"},
	{"/api/workspace/me/credentials", "GET"},
	{"/api/workspace/me/credentials/:wsc_id", "PUT"},
	{"/api/workspace/me/credentials/:wsc_id", "GET"},
	{"/api/workspace/me/credentials/:wsc_id/expiration", "PUT"},
}, PermissionOwnerSdk...)

var PermissionOwnerSdk = []authorizator.Permission{
	{"/api/sdk/application", "POST"},
	{"/api/sdk/application/:app_id", "PUT"},
	{"/api/sdk/application/:app_id", "DELETE"},
	{"/api/sdk/application", "GET"},
	{"/api/sdk/application/:app_id", "GET"},

	{"/api/sdk/application/:app_id/endpoint", "POST"},
	{"/api/sdk/application/:app_id/endpoint/:ep_id", "PUT"},
	{"/api/sdk/application/:app_id/endpoint/:ep_id", "DELETE"},
	{"/api/sdk/application/:app_id/endpoint", "GET"},
	{"/api/sdk/application/:app_id/endpoint/:ep_id", "GET"},

	{"/api/sdk/application/:app_id/endpoint/:ep_id/rule", "POST"},
	{"/api/sdk/application/:app_id/endpoint/:ep_id/rule/:epr_id", "PUT"},
	{"/api/sdk/application/:app_id/endpoint/:ep_id/rule/:epr_id", "DELETE"},
	{"/api/sdk/application/:app_id/endpoint/:ep_id/rule", "GET"},
	{"/api/sdk/application/:app_id/endpoint/:ep_id/rule/:epr_id", "GET"},

	{"/api/sdk/application/:app_id/message", "PUT"},
}
