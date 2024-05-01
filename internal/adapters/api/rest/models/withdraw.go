package models

import (
	"fmt"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Withdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func DeserializeWithdraw(data []byte) (*Withdraw, error) {
	withdraw := &Withdraw{}
	err := easyjson.Unmarshal(data, withdraw)
	if err != nil {
		return nil, fmt.Errorf("easyjson unmarshal error of deserialize withdraw:%w", err)
	}

	return withdraw, nil
}
