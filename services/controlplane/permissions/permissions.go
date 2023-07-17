package permissions

var (
	RoleAdmin = "admin"
)

var PermissionBase = [][]string{
	{"kanthor.controlplane.v1.Workspace", "Get"},
}

var PermissionAdmin = append(
	[][]string{
		{"kanthor.controlplane.v1.Workspace", "Update"},
	},
	PermissionBase...,
)
