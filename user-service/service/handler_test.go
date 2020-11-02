package service

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	database "keon.com/CitadelAllianceLobbyServer/user-service/dbclient"
	userModel "keon.com/CitadelAllianceLobbyServer/user-service/model"
	pb "keon.com/CitadelAllianceLobbyServer/user-service/proto/user"
)

func TestGetUser(t *testing.T) {
	tUser1 := &userModel.User{ID: "asdasd", Username: "test", HashedPassword: "123"}
	tUser2 := &userModel.User{}

	mockRepo := &database.MockRepository{}
	mockRepo.On("GetUserById", "asdasd").Return(tUser1, nil)
	mockRepo.On("GetUserById", "456").Return(tUser2, fmt.Errorf("Some error"))

	testCtx := context.Background()
	tokenService := &TokenService{mockRepo}
	ser := &service{repo: mockRepo, tokenService: tokenService}

	Convey("Given a grpc request to getUserbyid", t, func() {
		req := &pb.GetUserRequest{User: &pb.User{Id: "asdasd"}}
		res := &pb.GetUserResponse{}
		ser.GetUser(testCtx, req, res)

		Convey("Then the response should have no errors", func() {
			So(res.GetUser(), ShouldNotBeNil)

			us := res.GetUser()
			So(us.Id, ShouldEqual, "asdasd")
			So(us.Name, ShouldEqual, "test")
			So(us.Password, ShouldEqual, "123")
		})
	})
}

func TestGetAllUsers(t *testing.T) {
	tUser1 := &userModel.User{ID: "asdasd1", Username: "test1", HashedPassword: "123"}
	tUser2 := &userModel.User{ID: "asdasd2", Username: "test2", HashedPassword: "1234"}
	users := make([]*userModel.User, 0, 0)
	users = append(users, tUser1, tUser2)

	mockRepo := &database.MockRepository{}
	mockRepo.On("GetUsers").Return(users, nil)

	testCtx := context.Background()
	tokenService := &TokenService{mockRepo}
	ser := &service{repo: mockRepo, tokenService: tokenService}

	Convey("Given a grpc request to get all users", t, func() {
		req := &pb.GetAllRequest{}
		res := &pb.GetAllResponse{}
		ser.GetAllUsers(testCtx, req, res)

		Convey("Then the response should have no errors", func() {
			So(res.GetUsers(), ShouldNotBeNil)

			users := res.GetUsers()
			So(len(users), ShouldEqual, 2)

			u1 := users[0]
			u2 := users[1]

			So(u1.Id, ShouldEqual, "asdasd1")
			So(u1.Name, ShouldEqual, "test1")
			So(u1.Password, ShouldEqual, "123")

			So(u2.Id, ShouldEqual, "asdasd2")
			So(u2.Name, ShouldEqual, "test2")
			So(u2.Password, ShouldEqual, "1234")
		})
	})
}

func TestCreateUser(t *testing.T) {
	tUser1 := &userModel.User{ID: "asdasd1", Email: "asd@asd.com", Username: "test1", HashedPassword: "123"}
	tProtoUser1 := userModel.ToProtobufModel(tUser1)

	mockRepo := &database.MockRepository{}
	mockRepo.On("GetUserByEmail", "asd@asd.com").Return(nil, nil)
	mockRepo.On("GetUserByUsername", "test1").Return(nil, nil)
	mockRepo.On("CreateUser", tUser1).Return(nil)

	testCtx := context.Background()
	tokenService := &TokenService{mockRepo}
	ser := &service{repo: mockRepo, tokenService: tokenService}

	Convey("Given a grpc request to createUser", t, func() {
		req := &pb.CreateUserRequest{User: tProtoUser1}
		res := &pb.CreateUserResponse{}
		ser.CreateUser(testCtx, req, res)

		Convey("Then the response should have no errors", func() {
			So(res.Id, ShouldNotBeNil)

			//userID := res.Id

			//So(userID, ShouldEqual, "1")
		})
	})
}
