package utils

import (
	dto "gateway/DTO"
	"gateway/gen"

	"google.golang.org/protobuf/proto"
)

func UnmarshalListUsersReplay(data []byte) ([]dto.UserDTO, error) {
	protoUsers := &gen.ListUsersReplay{}

	err := proto.Unmarshal(data, protoUsers)
	if err != nil {
		return nil, err
	}

	userDTOs := []dto.UserDTO{}
	for _, protoUser := range protoUsers.Users {
		userDTO := dto.UserDTO{
			ID:        uint(protoUser.ID),
			FirstName: protoUser.FirstName,
			LastName:  protoUser.LastName,
			Email:     protoUser.Email,
		}

		userDTOs = append(userDTOs, userDTO)
	}

	return userDTOs, nil
}


func UnmarshalUser(data []byte) (*dto.UserDTO, error) {
	protoUser := &gen.User{}

	err := proto.Unmarshal(data, protoUser)
	if err != nil {
		return nil, err
	}

	userDTO := &dto.UserDTO{
		ID:        uint(protoUser.ID),
		FirstName: protoUser.FirstName,
		LastName:  protoUser.LastName,
		Email:     protoUser.Email,
	}

	return userDTO, nil
}

func MarshalCreateUserRequest(firstName, lastName, email string) ([]byte, error) {
	protoUser := &gen.CreateUserRequest{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	data, err := proto.Marshal(protoUser)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalUpdateUserRequest(userID uint, firstName, lastName string) ([]byte, error) {
	protoUser := &gen.UpdateUserRequest{
		ID:        uint32(userID),
		FirstName: firstName,
		LastName:  lastName,
	}

	data, err := proto.Marshal(protoUser)
	if err != nil {
		return nil, err
	}

	return data, nil
}
