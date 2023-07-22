package permissions

var Admin = "admin"

var AdminPermission = append(
	[][]string{
		{"kanthor.controlplane.v1.Workspace", "Update"},
	},
	BasePermission...,
)
