package permissions

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var Admin = "admin"

var AdminPermission = append(
	BasePermission,
	authorizator.Permission{Object: "kanthor.controlplane.v1.Workspace", Action: "Update"},
)
