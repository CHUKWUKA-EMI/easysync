package utils

import (
	"math/rand"
	"slices"
	"time"
)

// Generate8DigitsNumber ...
func Generate8DigitsNumber() int {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	return s.Intn(100000000)
}

// HasRequiredPermissions ...
func HasRequiredPermissions(currentRoles []interface{}, requiredRoles []string) bool {
	hasRequiredPermissions := false
	currentRoleNames := extractRoleNamesFromRoles(currentRoles)
	for _, role := range requiredRoles {
		if slices.Contains(currentRoleNames, role) {
			hasRequiredPermissions = true
			break
		}
	}
	return hasRequiredPermissions
}

func extractRoleNamesFromRoles(roles []interface{}) []string {
	var roleNames []string
	for _, role := range roles {
		r, ok := role.(map[string]any)
		if ok {
			roleName := r["name"].(string)
			roleNames = append(roleNames, roleName)
		}

	}
	return roleNames
}
