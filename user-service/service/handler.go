// vessel-service/handler.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/micro/go-micro/broker"
	"golang.org/x/crypto/bcrypt"
	userServiceDomain "keon.com/CitadelAllianceLobbyServer/user-service/dbclient/domain"
	userModel "keon.com/CitadelAllianceLobbyServer/user-service/model"
	pb "keon.com/CitadelAllianceLobbyServer/user-service/proto/user"
)

const CreateUserTopic = "user.created"

// Our gRPC Service handler
type Service struct {
	Repo         userServiceDomain.IUserRepository
	TokenService Authable
	PubSub       broker.Broker
}

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.Repo.GetUserById(ctx, req.User.GetId())
	if err != nil {
		return nil, err
	}
	us := user.(*userModel.User)
	res := &pb.GetUserResponse{User: userModel.ToProtobufModel(us)}

	return res, nil
}

func (s *Service) GetAllUsers(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := s.Repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := &pb.GetAllResponse{}
	sUsers := users.([]*userModel.User)
	for _, us := range sUsers {
		addUser := userModel.ToProtobufModel(us)
		res.Users = append(res.Users, addUser)
	}

	return res, nil
}

func (s *Service) Auth(ctx context.Context, req *pb.User) (*pb.Token, error) {
	log.Println("Logging in with:", req.Email, req.Password)
	_user, err := s.Repo.GetUserByEmail(ctx, req.Email)
	log.Println(_user)
	if err != nil {
		return nil, err
	}
	user := _user.(*userModel.User)
	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		return nil, err
	}

	token, err := s.TokenService.Encode(user)
	if err != nil {
		return nil, err
	}
	res := &pb.Token{Token: token, Errors: nil}

	return res, nil
}

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ref := req.User
	toCreate := &userModel.User{
		Username:       ref.GetName(),
		Email:          ref.GetEmail(),
		HashedPassword: ref.GetPassword(),
	}

	userByMail, err := s.Repo.GetUserByEmail(ctx, toCreate.Email)
	if err != nil {
		return nil, err
	}
	if userByMail != nil {
		return nil, errors.New("Email Already Exists")
	}

	userByUsername, err := s.Repo.GetUserByUsername(ctx, toCreate.Username)
	if err != nil {
		return nil, err
	}
	if userByUsername != nil {
		return nil, errors.New("Username Already Exists")
	}

	// Generates a hashed version of our password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(toCreate.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// req.User.Password = string(hashedPass)
	toCreate.HashedPassword = string(hashedPass)
	if err := s.Repo.CreateUser(ctx, toCreate); err != nil {
		return nil, err
	}

	// created, err := s.Repo.GetUserByEmail(ctx, toCreate.Email)
	res := &pb.CreateUserResponse{Id: toCreate.ID, Success: true}
	// if err := s.publishCreateUserResponse(res); err != nil {
	// 	return nil, err
	// }

	return res, nil
}

func (s *Service) ValidateToken(ctx context.Context, req *pb.Token) (*pb.Token, error) {

	// Decode token
	claims, err := s.TokenService.Decode(req.Token)
	if err != nil {
		return nil, err
	}

	log.Println(claims)

	if claims.User.ID == "" {
		return nil, errors.New("invalid user")
	}

	res := &pb.Token{Valid: true, Errors: nil}

	return res, nil
}

func (srv *Service) publishCreateUserResponse(res *pb.CreateUserResponse) error {
	// Marshal to JSON string
	body, err := json.Marshal(res)
	if err != nil {
		return err
	}

	// Create a broker message
	msg := &broker.Message{
		Header: map[string]string{
			"id":      res.Id,
			"success": strconv.FormatBool(res.Success),
		},
		Body: body,
	}

	// Publish message to broker
	if err := srv.PubSub.Publish(CreateUserTopic, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	}

	return nil
}
