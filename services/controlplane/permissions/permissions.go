package permissions

var (
	RoleOwner = "owner"
)

var (
	Policies = [][]string{
		{"kanthor.controlplane.v1.Account", "ListWorkspaces"},
	}
)

func PoliciesOfRoleInWorkspace(role, wsId string) [][]string {
	var policies [][]string

	for _, p := range Policies {
		policies = append(policies, append([]string{role, wsId}, p...))
	}

	return policies
}
