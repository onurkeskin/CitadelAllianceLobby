package model

import (
	pb "keon.com/CitadelAllianceLobbyServer/user-service-api/proto/user"
)

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

func FromProtobufModels(pbUsers []*pb.User) []*User {
	toReturn := make([]*User, 0)
	for _, pbUs := range pbUsers {
		toCreate := &User{
			Username:       pbUs.GetName(),
			Email:          pbUs.GetEmail(),
			HashedPassword: pbUs.GetPassword(),
		}
		toReturn = append(toReturn, toCreate)
	}

	return toReturn
}
