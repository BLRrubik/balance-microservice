package model

import "time"

type Accounting struct {
	Id        int        `json:"id"`
	UserId    int        `json:"userId"`
	Service   *Service   `json:"service"`
	OrderId   int        `json:"orderId"`
	Price     float64    `json:"price"`
	CreatedAt *time.Time `json:"created_at"`
}

type RevenueAddRequest struct {
	UserId    int     `json:"userId"`
	ServiceId int     `json:"serviceId"`
	OrderId   int     `json:"orderId"`
	Price     float64 `json:"price"`
}
