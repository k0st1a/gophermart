package rest

import "time"

type Register struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Withdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Order struct {
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
	Number     int64     `json:"number,string"`
	Accrual    float64   `json:"accrual,omitempty"`
}

type WithdrawOut struct {
	ProcessedAt time.Time `json:"processed_at"`
	Order       int64     `json:"order,string"`
	Sum         float64   `json:"sum"`
}
