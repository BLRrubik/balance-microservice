package sqlstore

import (
	"balance-microservice/internal/app/model"
	"balance-microservice/internal/app/store"
	"fmt"
	"strings"
)

type TransactionRepository struct {
	store *Store
}

func (repo *TransactionRepository) AddTransaction(t *model.Transaction) *store.Error {
	tx, err := repo.store.db.Begin()
	if err != nil {
		return &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	query := "INSERT INTO transactions (user_id, comment, price, status, date) " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING transaction_id"

	if err := tx.QueryRow(query, t.UserId, t.Comment, t.Price, t.Status, t.Date).Scan(&t.Id); err != nil {
		tx.Rollback()
		return &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	tx.Commit()

	return nil
}

func (repo *TransactionRepository) GetTransactionsOfUser(userId, page, size int,
	sort, dir string) ([]model.Transaction, *store.Error) {

	sort = strings.ToLower(sort)
	dir = strings.ToLower(dir)

	if sort != "date" && sort != "price" {
		return nil, &store.Error{
			Message:    "Invalid sorting param",
			StatusCode: 400,
		}
	}

	if dir != "asc" && dir != "desc" {
		return nil, &store.Error{
			Message:    "Invalid sort direction param",
			StatusCode: 400,
		}
	}

	if size <= 1 {
		return nil, &store.Error{
			Message:    "Size must be greater than 1",
			StatusCode: 400,
		}
	}

	if page < 0 {
		return nil, &store.Error{
			Message:    "Page must be positive",
			StatusCode: 400,
		}
	}

	if _, err := repo.store.UserRepository().FindById(userId); err != nil {
		return nil, err
	}

	transactions := make([]model.Transaction, 0)

	offset := page * size

	var query string

	if dir == "asc" {
		query = fmt.Sprintf("SELECT * FROM transactions WHERE user_id = $1 "+
			"ORDER BY %s OFFSET $2 LIMIT $3", sort)
	} else {
		query = fmt.Sprintf("SELECT * FROM transactions WHERE user_id = $1 "+
			"ORDER BY %s DESC OFFSET $2 LIMIT $3", sort)
	}

	rows, err := repo.store.db.Query(query, userId, offset, size)
	if err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	for rows.Next() {

		var t model.Transaction

		if err := rows.Scan(&t.Id, &t.UserId, &t.Price, &t.Comment, &t.Status, &t.Date); err != nil {
			return nil, &store.Error{
				Message:    err.Error(),
				StatusCode: 500,
			}
		}

		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return transactions, nil
}
