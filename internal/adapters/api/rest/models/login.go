package models

import (
	"fmt"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func DeserializeLogin(data []byte) (*Login, error) {
	login := &Login{}
	err := easyjson.Unmarshal(data, login)
	if err != nil {
		return nil, fmt.Errorf("easyjson unmarshal error of deserialize login:%w", err)
	}

	return login, nil
}
