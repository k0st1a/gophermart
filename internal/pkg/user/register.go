package user

import (
	"errors"
	"fmt"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Register struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var (
	ErrUserRegisterLoginEmpty    = errors.New("user register: login empty")
	ErrUserRegisterPasswordEmpty = errors.New("user register: password empty")
)

func DeserializeRegister(b []byte) (*Register, error) {
	r := newRegister()
	err := easyjson.Unmarshal(b, r)
	if err != nil {
		return nil, fmt.Errorf("user register unmarshal error:%w", err)
	}

	return r, nil
}

func newRegister() *Register {
	return &Register{}
}

func (r *Register) Validate() error {
	if r.Login == "" {
		return ErrUserRegisterLoginEmpty
	}

	if r.Password == "" {
		return ErrUserRegisterPasswordEmpty
	}

	return nil
}
