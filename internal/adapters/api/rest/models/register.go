package models

import (
	"fmt"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Register struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func DeserializeRegister(data []byte) (*Register, error) {
	register := &Register{}
	err := easyjson.Unmarshal(data, register)
	if err != nil {
		return nil, fmt.Errorf("easyjson unmarshal error of deserialize register:%w", err)
	}

	return register, nil
}
