package portalapi

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var (
	RoleOwner = "owner"
)

var PermissionOwner = []authorizator.Permission{
	{"/api/workspace/me", "GET"},
	{"/api/workspace/me", "PUT"},

	{"/api/workspace/me/credentials", "POST"},
	{"/api/workspace/me/credentials", "GET"},
	{"/api/workspace/me/credentials/:wsc_id", "PUT"},
	{"/api/workspace/me/credentials/:wsc_id", "GET"},
	{"/api/workspace/me/credentials/:wsc_id/expiration", "PUT"},
}
