package permissions

import "github.com/scrapnode/kanthor/infrastructure/authorizator"

var BasePermission = []authorizator.Permission{
	{"kanthor.dataplane.v1.Msg", "Put"},
}
