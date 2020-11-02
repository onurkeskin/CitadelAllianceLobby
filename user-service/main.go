package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"

	_ "github.com/micro/go-plugins/broker/nats"
	"google.golang.org/grpc/credentials"
	dbclient "keon.com/CitadelAllianceLobbyServer/user-service/dbclient"
	pb "keon.com/CitadelAllianceLobbyServer/user-service/proto/user"
	service "keon.com/CitadelAllianceLobbyServer/user-service/service"
)

const (
	defaultHost = "localhost:27017"
)

func main() {

	// Creates a database connection and handles
	// closing it again before exit.
	db, err := dbclient.CreateConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	//FRONT-END CERTIFICATES
	// frontendCert, err := ioutil.ReadFile("certs/frontend.crt")
	// if err != nil {
	// 	log.Fatalf("Could not read front end cert file: %v", err)
	// }
	// rootsFront := x509.NewCertPool()
	// rootsFront.AppendCertsFromPEM(frontendCert)
	// credsClient := credentials.NewClientTLSFromCert(rootsFront, "")

	//BACK-END CERTIFICATES
	backendCert, err := ioutil.ReadFile("certs/backend.cert")
	if err != nil {
		log.Fatalf("Could not read backend cert file: %v", err)
	}
	backendKey, err := ioutil.ReadFile("certs/backend.key")
	if err != nil {
		log.Fatalf("Could not read backend key file: %v", err)
	}
	certBackend, err := tls.X509KeyPair(backendCert, backendKey)
	if err != nil {
		log.Fatalf("failed to parse certificate: %v", err)
	}
	creds := credentials.NewServerTLSFromCert(&certBackend)

	serverOption := grpc.Creds(creds)
	server := grpc.NewServer(serverOption)
	defer server.Stop()

	repo := &dbclient.UserRepository{db}

	tokenService := &service.TokenService{Repo: repo}

	// pubsub := srv.Server().Options().Broker
	// if err := pubsub.Connect(); err != nil {
	// 	fmt.Print("Broker not online")
	// 	log.Fatal(err)
	// }

	// Register handler
	pb.RegisterUserServiceServer(server, &service.Service{Repo: repo, TokenService: tokenService})

	lis, err := net.Listen("tcp", ":6767")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Starting to list on %d", 8080)
	// Run the server
	if err := server.Serve(lis); err != nil {
		log.Println(err)
	}
}
