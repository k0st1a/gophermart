package accrual

import (
	"fmt"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func Deserialize(data []byte) (*Accrual, error) {
	accrual := &Accrual{}
	err := easyjson.Unmarshal(data, accrual)
	if err != nil {
		return nil, fmt.Errorf("easyjson unmarshal error of deserialize order:%w", err)
	}

	return accrual, nil
}
