package service

import (
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	domain "keon.com/CitadelAllianceLobbyServer/user-service-api/domain"
	serviceDomain "keon.com/CitadelAllianceLobbyServer/user-service-api/service/domain"
)

type PostCreateUserHookPayload struct {
	User domain.IUser
}

type PostConfirmUserHookPayload struct {
	User domain.IUser
}

type ControllerHooks struct {
	PostCreateUserHook  func(resource *Resource, w http.ResponseWriter, req *http.Request, payload *PostCreateUserHookPayload) error
	PostConfirmUserHook func(resource *Resource, w http.ResponseWriter, req *http.Request, payload *PostConfirmUserHookPayload) error
}

type Options struct {
	BasePath        string
	Renderer        serviceDomain.IRenderer
	ControllerHooks *ControllerHooks
	GrpcHelper      *GrpcHelper
}

type GrpcHelper struct {
	GrpcAddress    string
	GrpcClientCert *credentials.TransportCredentials
}

func (g *GrpcHelper) getGRPCClient() (*grpc.ClientConn, error) {
	// Dial with specific Transport (with credentials)
	conn, err := grpc.Dial(g.GrpcAddress, grpc.WithTransportCredentials(*g.GrpcClientCert))
	if err != nil {
		return nil, nil
	}

	return conn, nil
}

func NewResource(options *Options) *Resource {
	renderer := options.Renderer
	if renderer == nil {
		panic("users.Options.Renderer is required")
	}

	controllerHooks := options.ControllerHooks
	if controllerHooks == nil {
		controllerHooks = &ControllerHooks{nil, nil}
	}

	u := &Resource{
		options,
		nil,
		renderer,
		controllerHooks,
		options.GrpcHelper,
	}
	u.generateRoutes(options.BasePath)
	return u
}

// UsersResource implements IResource
type Resource struct {
	options         *Options
	routes          *serviceDomain.Routes
	Renderer        serviceDomain.IRenderer
	ControllerHooks *ControllerHooks
	GrpcHelper      *GrpcHelper
}

func (resource *Resource) Routes() *serviceDomain.Routes {
	return resource.routes
}

// func (resource *Resource) GetClient() pb.UserServiceClient {
// 	return resource.pbClientCert
// }

func (resource *Resource) Render(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
	resource.Renderer.Render(w, req, status, v)
}
