package user

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/mailru/easyjson"
)

//easyjson:json
type Register struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

var (
	ErrUserRegisterLoginEmpty    = errors.New("user register: login empty")
	ErrUserRegisterPasswordEmpty = errors.New("user register: password empty")
	ErrUserRegisterValidation    = errors.New("user register: validation error")
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
	v := validator.New()
	err := v.Struct(*r)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrUserRegisterValidation, err)
	}

	if r.Login == "" {
		return ErrUserRegisterLoginEmpty
	}

	if r.Password == "" {
		return ErrUserRegisterPasswordEmpty
	}

	return nil
}
