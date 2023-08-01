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
}
