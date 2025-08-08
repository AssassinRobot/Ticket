package utils

import (
	"errors"
	"strings"
)

func HandleError(data []byte)error{
	strData := string(data)

	if !strings.Contains(strData,"error"){
		return  nil
	}

	return  errors.New(strData)
}