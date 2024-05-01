package models

import (
	"fmt"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func SerializeBalance(b *Balance) ([]byte, error) {
	data, err := easyjson.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("easyjson marshal error of serialize balance:%w", err)
	}

	return data, nil
}
