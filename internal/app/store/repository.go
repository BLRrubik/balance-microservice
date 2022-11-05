package store

import (
	"balance-microservice/internal/app/model"
)

type UserRepository interface {
	ReplenishmentAccount(request *model.UserDepositRequest) (*model.User, *Error)
	GetBalanceOfUser(int) (*model.User, *Error)
	FindAll() ([]model.User, *Error)
	FindById(id int) (*model.User, *Error)
	UpdateBalance(u *model.User, price float64) (*model.User, *Error)
}

type ServiceRepository interface {
	FindById(id int) (*model.Service, *Error)
}

type BillRepository interface {
	ReservedFunds(request *model.BillCreateRequest) (*model.Bill, *Error)
	ApproveReservation(id int) (*model.Bill, *Error)
	RejectReservation(id int) (*model.Bill, *Error)
	FindById(id int) (*model.Bill, *Error)
	FindByIdAndStatus(id int, status model.Status) (*model.Bill, *Error)
}

type AccountingRepository interface {
	AddRevenue(request *model.RevenueAddRequest) *Error
	FindAll() ([]model.Accounting, *Error)
	ExportCSV(dt string) (*string, *Error)
}

type TransactionRepository interface {
	AddTransaction(transaction *model.Transaction) *Error
	GetTransactionsOfUser(userId, page, size int, sort, dir string) ([]model.Transaction, *Error)
}
