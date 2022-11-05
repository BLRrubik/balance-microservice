package sqlstore

import (
	"balance-microservice/internal/app/model"
	"balance-microservice/internal/app/store"
	"database/sql"
	"fmt"
	"time"
)

type BillRepository struct {
	store *Store
}

func (repo *BillRepository) ApproveReservation(id int) (*model.Bill, *store.Error) {

	bill, e := repo.store.BillRepository().FindByIdAndStatus(id, model.PENDING)

	if e != nil {
		e.Message = "Record is not found or not at 'PENDING' status"
		return nil, e
	}

	tx, err := repo.store.db.Begin()
	if err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	query := "UPDATE bills SET status = '" + model.APPROVED + "', moderate_at = now() " +
		"WHERE bill_id = $1 RETURNING bill_id, status, moderate_at"

	if err := tx.QueryRow(query, bill.Id).Scan(&bill.Id, &bill.Status, &bill.ModerateAt); err != nil {
		tx.Rollback()
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	req := &model.RevenueAddRequest{
		UserId:    bill.User.Id,
		ServiceId: bill.Service.Id,
		OrderId:   bill.OrderId,
		Price:     bill.Price,
	}

	if err := repo.store.AccountingRepository().AddRevenue(req); err != nil {
		return nil, err
	}

	transaction := &model.Transaction{
		UserId:  bill.User.Id,
		Comment: fmt.Sprintf("Approve transaction for service: '%s'", bill.Service.Name),
		Date:    time.Now(),
		Status:  model.PAYMENT,
		Price:   bill.Price,
	}

	if err := repo.store.TransactionRepository().AddTransaction(transaction); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return bill, nil
}

func (repo *BillRepository) RejectReservation(id int) (*model.Bill, *store.Error) {

	bill, e := repo.store.BillRepository().FindByIdAndStatus(id, model.PENDING)
	if e != nil {
		e.Message = "Record is not found or not at 'PENDING' status"
		return nil, e
	}

	tx, err := repo.store.db.Begin()
	if err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	query := "UPDATE bills SET status = '" + model.REJECT + "', moderate_at = now() " +
		"WHERE bill_id = $1 RETURNING bill_id, status, moderate_at"

	if err := tx.QueryRow(query, bill.Id).Scan(&bill.Id, &bill.Status, &bill.ModerateAt); err != nil {
		tx.Rollback()
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	if _, err := repo.store.UserRepository().UpdateBalance(bill.User, bill.User.Balance+bill.Price); err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction := &model.Transaction{
		UserId:  bill.User.Id,
		Comment: fmt.Sprintf("Reject transaction for service: '%s'", bill.Service.Name),
		Date:    time.Now(),
		Status:  model.DEPOSIT,
		Price:   bill.Price,
	}

	if err := repo.store.TransactionRepository().AddTransaction(transaction); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return bill, nil
}

func (repo *BillRepository) ReservedFunds(request *model.BillCreateRequest) (*model.Bill, *store.Error) {

	tx, err := repo.store.db.Begin()
	if err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	user, e := repo.store.UserRepository().FindById(request.UserId)
	if e != nil {
		return nil, e
	}

	service, e := repo.store.ServiceRepository().FindById(request.ServiceId)
	if e != nil {
		return nil, e
	}

	if request.Price <= 0 {
		return nil, &store.Error{
			Message:    "Price must be greater than zero",
			StatusCode: 400,
		}
	}

	if request.OrderId <= 0 {
		return nil, &store.Error{
			Message:    "OrderId must be greater than zero",
			StatusCode: 400,
		}
	}

	balance := user.Balance - request.Price

	if balance < 0 {
		return nil, &store.Error{
			Message:    "Not enough money on balance",
			StatusCode: 400,
		}
	}

	bill := &model.Bill{
		User:    user,
		Service: service,
		OrderId: request.OrderId,
		Price:   request.Price,
		Status:  model.PENDING,
	}

	query := "INSERT INTO bills (user_id, service_id, order_id, price) " +
		"VALUES ($1, $2, $3, $4) RETURNING bill_id"

	if err := tx.QueryRow(query, user.Id, service.Id, request.OrderId, request.Price).
		Scan(&bill.Id); err != nil {
		tx.Rollback()
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	if _, err := repo.store.UserRepository().UpdateBalance(user, balance); err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction := &model.Transaction{
		UserId:  bill.User.Id,
		Comment: fmt.Sprintf("Reserved funds for service: '%s'", bill.Service.Name),
		Date:    time.Now(),
		Status:  model.RESERVED,
		Price:   bill.Price,
	}

	if err := repo.store.TransactionRepository().AddTransaction(transaction); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return bill, nil
}

func (repo *BillRepository) FindById(id int) (*model.Bill, *store.Error) {
	query := "SELECT * FROM bill WHERE bill_id = $1"

	b := &model.Bill{}
	u := &model.User{}
	s := &model.Service{}

	if err := repo.store.db.QueryRow(
		query,
		id,
	).Scan(&b.Id, &b.OrderId, &b.Price, &b.Status, &b.ModerateAt,
		&u.Id, &s.Id, &s.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, &store.Error{
				Message:    fmt.Sprintf("Bill with id %d is not found", id),
				StatusCode: 404,
			}
		}

		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	b.Service = s
	b.User = u

	return b, nil
}

func (repo *BillRepository) FindByIdAndStatus(id int, status model.Status) (*model.Bill, *store.Error) {

	query := "SELECT * FROM bill WHERE bill_id = $1 AND status = $2"

	b := &model.Bill{}
	u := &model.User{}
	s := &model.Service{}

	if err := repo.store.db.QueryRow(
		query,
		id,
		status,
	).Scan(&b.Id, &b.OrderId, &b.Price, &b.Status, &b.ModerateAt,
		&u.Id, &u.Balance, &s.Id, &s.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, &store.Error{
				Message:    fmt.Sprintf("Bill with id %d is not found", id),
				StatusCode: 404,
			}
		}

		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	b.Service = s
	b.User = u

	return b, nil

}
