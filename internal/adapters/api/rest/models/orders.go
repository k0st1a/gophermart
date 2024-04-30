package models

import (
	"fmt"
	"time"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Order struct {
	Number     int64     `json:"number,string"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
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
