package models

import (
	"fmt"
	"time"

	"github.com/mailru/easyjson"
)

//easyjson:json
type WithdrawOut struct {
	ProcessedAt time.Time `json:"processed_at"`
	Order       int64     `json:"order,string"`
	Sum         float64   `json:"sum"`
}

//easyjson:json
type Withdrawals []WithdrawOut

func SerializeWithdrawals(w *Withdrawals) ([]byte, error) {
	data, err := easyjson.Marshal(w)
	if err != nil {
		return nil, fmt.Errorf("easyjson marshal error of serialize withdrawals:%w", err)
	}

	return data, nil
}
