package model

import (
	"github.com/volatiletech/null"
	caModels "keon.com/CitadelAllianceLobbyServer/user-service/dbmodels"
	pb "keon.com/CitadelAllianceLobbyServer/user-service/proto/user"
)

func FromDBModel(us *caModels.User) *User {
	modelUser := &User{
		ID:               us.ID,
		Username:         us.Username.String,
		Email:            us.Email.String,
		HashedPassword:   us.R.UserPassword.HashedPassword.String,
		Status:           us.Status.String,
		ConfirmationCode: us.ConfirmationCode.String,
		LastModifiedDate: us.ModifiedDate.Time,
	}

	for _, r := range us.Roles {
		modelUser.Roles = append(modelUser.Roles, r)
	}

	return modelUser
}

func ToDBModel(us *User) *caModels.User {
	var toCreateUser = &caModels.User{
		ID:       us.ID,
		Username: null.StringFrom(us.Username),
		Email:    null.StringFrom(us.Email),
		Status:   null.StringFrom("Created"),
	}

	for _, str := range us.Roles {
		toCreateUser.Roles = append(toCreateUser.Roles, str)
	}

	return toCreateUser
}

func FromProtobufModel(us *pb.User) *User {
	toCreate := &User{
		Username:       us.GetName(),
		Email:          us.GetEmail(),
		HashedPassword: us.GetPassword(),
	}
	return toCreate
}

func ToProtobufModel(us *User) *pb.User {
	toReturn := &pb.User{}
	if us.ID != "" {
		toReturn.Id = us.ID
	}

	if us.Email != "" {
		toReturn.Email = us.Email
	}

	if us.Username != "" {
		toReturn.Name = us.Username
	}

	if us.HashedPassword != "" {
		toReturn.Password = us.HashedPassword
	}

	return toReturn
}
