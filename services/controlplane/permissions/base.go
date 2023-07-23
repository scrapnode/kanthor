package permissions

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var BasePermission = []authorizator.Permission{
	{"kanthor.controlplane.v1.Workspace", "Get"},
}
