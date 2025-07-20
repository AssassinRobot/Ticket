package utils

import (
	"user/gen"
	"user/internal/application/core/domain"

	"google.golang.org/protobuf/proto"
)

func Marshal(user *domain.User) ([]byte,error) {
	protoUser := &gen.User{
		ID:        uint32(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	data, err := proto.Marshal(protoUser)
	if err != nil {
		return nil, err
	}

	return data,nil

}
