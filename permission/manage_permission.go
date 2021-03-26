package permission

import (
	"strings"
)

type permission struct {
	roles   []string
	methods []string
}

type authority map[string]permission

var authorities = authority{
	"/user": permission{
		roles:		[]string{"USER"},
		methods:	[]string{"GET", "PUT", "DELETE"},
	},
	"/api": permission{
		roles:		[]string{"USER"},
		methods:	[]string{"GET"},
	},
	"/admin": permission{
		roles:		[]string{"ADMIN"},
		methods:	[]string{"GET", "POST", "PUT", "DELETE"},
	},
	"/login": permission{
		roles:   []string{"USER", "ADMIN"},
		methods: []string{"POST"},
	},
	"/logout": permission{
		roles:   []string{"USER", "ADMIN"},
		methods: []string{"POST"},
	},
	"/signup": permission{
		roles:   []string{"USER"},
		methods: []string{"POST"},
	},
}

// HasPermission checks if a given role has permission to access a given route for a given method
func HasPermission(path string, role string, method string) bool {
	if strings.HasPrefix(path, "/admin") {
		path = "/admin"
	}
	perm := authorities[path]
	checkedRole := checkRole(role, perm.roles)
	checkedMethod := checkMethod(method, perm.methods)
	if !checkedRole || !checkedMethod {
		return false
	}
	return true
}

func checkRole(role string, roles []string) bool {
	for _, r := range roles {
		if strings.ToUpper(r) == strings.ToUpper(role) {
			return true
		}
	}
	return false
}

func checkMethod(method string, methods []string) bool {
	for _, m := range methods {
		if strings.ToUpper(m) == strings.ToUpper(method) {
			return true
		}
	}
	return false
}