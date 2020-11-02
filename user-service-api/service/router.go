package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"keon.com/CitadelAllianceLobbyServer/user-service-api/service/domain"
)

// Router type
type Router struct {
	*mux.Router
}

// matcherFunc matches the handler to the correct API version based on its `accept` header
// TODO: refactor matcher function as server.Config
func matcherFunc(r domain.Route, defaultHandler http.HandlerFunc) func(r *http.Request, rm *mux.RouteMatch) bool {
	return func(req *http.Request, rm *mux.RouteMatch) bool {
		acceptHeaders := domain.NewAcceptHeadersFromString(req.Header.Get("accept"))
		foundHandler := defaultHandler
		// try to match a handler to the specified `version` params
		// else we will fall back to the default handler
		for _, h := range acceptHeaders {
			m := h.MediaType
			// check if media type is `application/json` type or `application/[*]+json` suffix
			if !(m.Type == "application" && (m.SubType == "json" || m.Suffix == "json")) {
				continue
			}

			// if its the right application type, check if a version specified
			version, hasVersion := m.Parameters["version"]
			if !hasVersion {
				continue
			}
			if handler, ok := r.RouteHandlers[domain.RouteHandlerVersion(version)]; ok {
				// found handler for specified version
				foundHandler = handler
				break
			}
		}

		// if ac != nil {
		// 	rm.Handler = ac.NewContextHandler(r.Name, foundHandler)
		// } else {
		// 	rm.Handler = foundHandler
		// }

		rm.Handler = foundHandler
		return true
	}
}

// NewRouter Returns a new Router object
func NewRouter() *Router {
	router := mux.NewRouter().StrictSlash(true)

	return &Router{router}
}

func (router *Router) AddRoutes(routes *domain.Routes) *Router {
	if routes == nil {
		return router
	}
	for _, route := range *routes {

		// get the defaultHandler for current route at init time so that we can safely panic
		// if it was not defined
		defaultHandler, ok := route.RouteHandlers[route.DefaultVersion]
		if !ok {
			// server/router instantiation error
			// its safe to throw panic here
			panic(errors.New(fmt.Sprintf("Routes definition error, missing default route handler for version `%v` in `%v`",
				route.DefaultVersion, route.Name)))
		}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Queries(route.Queries...).
			MatcherFunc(matcherFunc(route, defaultHandler))
	}
	return router
}

func (router *Router) AddResources(resources ...domain.IResource) *Router {
	for _, resource := range resources {
		if resource.Routes() == nil {
			// server/router instantiation error
			// its safe to throw panic here
			panic(errors.New(fmt.Sprintf("Routes definition missing: %v", resource)))
		}
		router.AddRoutes(resource.Routes())
	}
	return router
}
