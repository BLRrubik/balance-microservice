package model

import "time"

type Status string

const (
	PENDING  Status = "PENDING"
	REJECT          = "REJECT"
	APPROVED        = "APPROVED"
)

type Bill struct {
	Id         int
	User       *User
	Service    *Service
	OrderId    int
	Price      float64
	Status     Status
	ModerateAt *time.Time
}

type BillCreateRequest struct {
	UserId    int     `json:"userId"`
	ServiceId int     `json:"serviceId"`
	OrderId   int     `json:"orderId"`
	Price     float64 `json:"price"`
}

type BillDto struct {
	Id      int `json:"id"`
	UserId  int `json:"userId"`
	Service struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"service"`
	OrderId    int        `json:"orderId"`
	Price      float64    `json:"price"`
	Status     Status     `json:"status"`
	ModerateAt *time.Time `json:"moderateAt,omitempty"`
}

func (o *Bill) ToDto() *BillDto {
	type service struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	return &BillDto{
		Id:     o.Id,
		UserId: o.User.Id,
		Service: service{
			Id:   o.Service.Id,
			Name: o.Service.Name,
		},
		OrderId:    o.OrderId,
		Price:      o.Price,
		Status:     o.Status,
		ModerateAt: o.ModerateAt,
	}
}
