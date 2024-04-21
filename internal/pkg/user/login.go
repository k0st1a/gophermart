package user

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/mailru/easyjson"
)

//easyjson:json
type Login struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

var (
	ErrUserLoginLoginEmpty    = errors.New("user login: login empty")
	ErrUserLoginPasswordEmpty = errors.New("user login: password empty")
	ErrUserLoginValidation    = errors.New("user login: validation error")
)

func DeserializeLogin(b []byte) (*Login, error) {
	r := newLogin()
	err := easyjson.Unmarshal(b, r)
	if err != nil {
		return nil, fmt.Errorf("user login unmarshal error:%w", err)
	}

	return r, nil
}

func newLogin() *Login {
	return &Login{}
}

func (r *Login) Validate() error {
	v := validator.New()
	err := v.Struct(*r)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrUserLoginValidation, err)
	}

	if r.Login == "" {
		return ErrUserLoginLoginEmpty
	}

	if r.Password == "" {
		return ErrUserLoginPasswordEmpty
	}

	return nil
}
