package main

import (
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/kr/pretty"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/credentials"

	service "keon.com/CitadelAllianceLobbyServer/user-service-api/service"
	renderer "keon.com/CitadelAllianceLobbyServer/user-service-api/service/middlewares/renderer"
)

var appName = "user-service-api"

func main() {
	// Read cert file
	FrontendCert, _ := ioutil.ReadFile("./certs/frontend.cert")

	// Create CertPool
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(FrontendCert)

	// Create credentials
	credsClient := credentials.NewClientTLSFromCert(roots, "")

	renderer := renderer.New(&renderer.Options{
		IndentJSON: true,
	}, renderer.JSON)

	controllerHooks := service.ControllerHooks{
		PostCreateUserHook:  nil,
		PostConfirmUserHook: nil,
	}

	usersResource := service.NewResource(&service.Options{
		Renderer:        renderer,
		ControllerHooks: &controllerHooks,
		GrpcHelper:      &service.GrpcHelper{GrpcAddress: "user-service:6767", GrpcClientCert: &credsClient},
	})

	log.Printf("Starting %v\n", appName)

	// init server
	s := service.NewServer(&service.Config{})
	router := service.NewRouter()

	// add REST resources to router
	router.AddResources(usersResource)

	// setup router
	s.UseRouter(router)

	KeyPath := "certs/user-api-private.key"
	CertPath := "certs/user-api-public.crt"

	f, err := os.OpenFile("./logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("asdasd")

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		log.Println(pretty.Sprint(route))
		return nil
	})

	s.Run(":6767", service.ServerOptions{
		Timeout:  10 * time.Second,
		CertPath: CertPath,
		KeyPath:  KeyPath,
	})
}
