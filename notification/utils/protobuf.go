package utils

import (
	"fmt"
	"notification/gen"

	"google.golang.org/protobuf/proto"
)

func UnMarshalUserData(data []byte) (string, string, error) {
	protoUser := &gen.User{}

	err := proto.Unmarshal(data, protoUser)
	if err != nil {
		return "", "", err
	}

	fullname := fmt.Sprintf("%s %s", protoUser.FirstName, protoUser.LastName)

	return fullname, protoUser.Email, nil

}
