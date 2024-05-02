package models

import (
	"fmt"
	"time"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Order struct {
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
	Number     int64     `json:"number,string"`
	Accrual    float64   `json:"accrual,omitempty"`
}

//easyjson:json
type Orders []Order

func SerializeOrders(o *Orders) ([]byte, error) {
	data, err := easyjson.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("easyjson marshal error of serialize orders:%w", err)
	}

	return data, nil
}
