package authorizator

func PoliciesOfRoleInWorkspace(role, wsId string, permissions [][]string) [][]string {
	var policies [][]string

	for _, p := range permissions {
		policies = append(policies, append([]string{role, wsId}, p...))
	}

	return policies
}
