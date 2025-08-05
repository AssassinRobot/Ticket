package utils

import (
	"user/gen"
	"user/internal/application/core/domain"

	"google.golang.org/protobuf/proto"
)

func UnmarshalCreateUserRequest(data []byte) (string, string, string, error) {
	protoUser := &gen.CreateUserRequest{}

	err := proto.Unmarshal(data, protoUser)
	if err != nil {
		return "", "", "", err
	}

	return protoUser.FirstName, protoUser.LastName, protoUser.Email, nil
}

func UnmarshalUpdateUserRequest(data []byte) (uint, string, string, error) {
	protoUser := &gen.UpdateUserRequest{}

	err := proto.Unmarshal(data, protoUser)
	if err != nil {
		return 0, "", "", err
	}

	return uint(protoUser.ID),protoUser.FirstName, protoUser.LastName, nil
}

func MarshalUser(user *domain.User) ([]byte, error) {
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

	return data, nil
}

func MarshalListUsersReplay(users []domain.User) ([]byte, error) {
	protoListUserReplay := &gen.ListUsersReplay{}

	for _, user := range users {
		protoUser := &gen.User{
			ID:        uint32(user.ID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}

		protoListUserReplay.Users = append(protoListUserReplay.Users, protoUser)
	}

	data, err := proto.Marshal(protoListUserReplay)
	if err != nil {
		return nil, err
	}

	return data, nil
}

