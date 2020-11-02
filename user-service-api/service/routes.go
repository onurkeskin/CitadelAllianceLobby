package service

import (
	"strings"

	serviceDomain "keon.com/CitadelAllianceLobbyServer/user-service-api/service/domain"
)

const (
	ListUsers = "ListUsers"
	// CountUsers       = "CountUsers"
	// GetUser          = "GetUser"
	// CreateUser       = "CreateUser"
	// UpdateUsers      = "UpdateUsers"
	// DeleteAllUsers   = "DeleteAllUsers"
	// ConfirmUser      = "ConfirmUser"
	// UpdateUser       = "UpdateUser"
	// DeleteUser       = "DeleteUser"
	// ConfirmAgainUser = "ResendConfirmationUser"
)

const defaultBasePath = "/users-api"

func (resource *Resource) generateRoutes(basePath string) *serviceDomain.Routes {
	if basePath == "" {
		basePath = defaultBasePath
	}
	var baseRoutes = serviceDomain.Routes{
		serviceDomain.Route{
			Name:           ListUsers,
			Method:         "GET",
			Pattern:        "/users-api/users",
			DefaultVersion: "0.0",
			RouteHandlers: serviceDomain.RouteHandlers{
				"0.0": resource.HandleListUsers_v0,
			},
			// ACLHandler: resource.HandleListUsersACL,
		},
	}

	routes := serviceDomain.Routes{}

	for _, route := range baseRoutes {
		r := serviceDomain.Route{
			Name:           route.Name,
			Method:         route.Method,
			Pattern:        strings.Replace(route.Pattern, defaultBasePath, basePath, -1),
			DefaultVersion: route.DefaultVersion,
			RouteHandlers:  route.RouteHandlers,
			// ACLHandler:     route.ACLHandler,
		}
		routes = routes.Append(&serviceDomain.Routes{r})
	}
	resource.routes = &routes
	return resource.routes
}
