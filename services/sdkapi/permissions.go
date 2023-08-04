package sdkapi

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	RoleOwner = "owner"
)

var PermissionOwner = []authorizator.Permission{
	{"/application", "POST"},
	{"/application/:app_id", "PATCH"},
	{"/application/:app_id", "DELETE"},
	{"/application", "GET"},
	{"/application/:app_id", "GET"},

	{"/application/:app_id/endpoint", "POST"},
	{"/application/:app_id/endpoint/:ep_id", "PATCH"},
	{"/application/:app_id/endpoint/:ep_id", "DELETE"},
	{"/application/:app_id/endpoint", "GET"},
	{"/application/:app_id/endpoint/:ep_id", "GET"},

	{"/application/:app_id/endpoint/:ep_id/rule", "POST"},
	{"/application/:app_id/endpoint/:ep_id/rule/:epr_id", "PATCH"},
	{"/application/:app_id/endpoint/:ep_id/rule/:epr_id", "DELETE"},
	{"/application/:app_id/endpoint/:ep_id/rule", "GET"},
	{"/application/:app_id/endpoint/:ep_id/rule/:epr_id", "GET"},
}
