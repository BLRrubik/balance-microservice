package model

import "time"

type TransactionStatus string

const (
	DEPOSIT  TransactionStatus = "DEPOSIT"
	PAYMENT                    = "PAYMENT"
	RESERVED                   = "RESERVED"
)

type Transaction struct {
	Id      int               `json:"id"`
	UserId  int               `json:"userId"`
	Price   float64           `json:"price"`
	Comment string            `json:"comment"`
	Status  TransactionStatus `json:"status"`
	Date    time.Time         `json:"date"`
}
