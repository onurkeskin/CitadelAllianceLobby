package main

import (
	"log"
	"os"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"golang.org/x/net/context"
	pb "keon.com/CitadelAllianceLobbyServer/user-service/proto/user"
)

func main() {

	cmd.Init()

	// Create new greeter client
	client := pb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)

	name := "Onur Keskin"
	email := "onurkeskindev@gmail.com"
	password := "test123"

	log.Println(name, email, password)

	r, err := client.CreateUser(context.TODO(), &pb.CreateUserRequest{
		User: &pb.User{
			Name:     name,
			Email:    email,
			Password: password,
		},
	})
	if err != nil {
		log.Printf("Could not create: %v", err)
	} else {
		log.Printf("Created: %s", r.Success)
	}

	getAll, err := client.GetAllUsers(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatalf("Could not list users: %v", err)
	}
	for _, v := range getAll.Users {
		log.Println(v)
	}

	authResponse, err := client.Auth(context.TODO(), &pb.User{
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Could not authenticate user: %s error: %v\n", email, err)
	}

	log.Printf("Your access token is: %s \n", authResponse.Token)

	// let's just exit because
	os.Exit(0)
}
