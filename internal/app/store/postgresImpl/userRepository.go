package sqlstore

import (
	"balance-microservice/internal/app/model"
	"balance-microservice/internal/app/store"
	"database/sql"
	"fmt"
	"time"
)

type UserRepository struct {
	store *Store
}

func (repo *UserRepository) ReplenishmentAccount(req *model.UserDepositRequest) (*model.User, *store.Error) {
	user, err := repo.FindById(req.Id)
	transaction := &model.Transaction{}

	if err != nil {

		tx, err := repo.store.db.Begin()
		if err != nil {
			return nil, &store.Error{
				Message:    err.Error(),
				StatusCode: 500,
			}
		}

		query := "INSERT INTO users (user_id, balance) VALUES ($1, $2) RETURNING user_id, balance"

		user := model.User{}

		if err := repo.store.db.QueryRow(query, req.Id, req.Deposit).Scan(&user.Id, &user.Balance); err != nil {
			tx.Rollback()
			return nil, &store.Error{
				Message:    err.Error(),
				StatusCode: 500,
			}
		}

		transaction.UserId = user.Id
		transaction.Date = time.Now()
		transaction.Comment = "Deposit on balance"
		transaction.Status = model.DEPOSIT
		transaction.Price = req.Deposit

		if err := repo.store.TransactionRepository().AddTransaction(transaction); err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Commit()

		return &user, nil
	}

	user.Balance += req.Deposit

	if _, err := repo.UpdateBalance(user, user.Balance); err != nil {
		return nil, err
	}

	transaction.UserId = user.Id
	transaction.Date = time.Now()
	transaction.Comment = "Deposit on balance"
	transaction.Status = model.DEPOSIT
	transaction.Price = req.Deposit

	if err := repo.store.TransactionRepository().AddTransaction(transaction); err != nil {
		return nil, err
	}

	return user, nil

}

func (repo *UserRepository) GetBalanceOfUser(id int) (*model.User, *store.Error) {
	u, err := repo.FindById(id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (repo *UserRepository) FindById(id int) (*model.User, *store.Error) {
	query := "SELECT * FROM users WHERE user_id = $1"

	u := &model.User{}
	if err := repo.store.db.QueryRow(
		query,
		id,
	).Scan(&u.Id, &u.Balance); err != nil {
		if err == sql.ErrNoRows {
			return nil, &store.Error{
				Message:    fmt.Sprintf("User with id %d is not found", id),
				StatusCode: 404,
			}
		}

		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return u, nil
}

func (repo *UserRepository) FindAll() ([]model.User, *store.Error) {
	query := "SELECT * FROM users"

	users := make([]model.User, 0)

	rows, err := repo.store.db.Query(query)
	if err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Balance); err != nil {
			return nil, &store.Error{
				Message:    err.Error(),
				StatusCode: 500,
			}
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return users, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return users, nil
}

func (repo *UserRepository) UpdateBalance(u *model.User, balance float64) (*model.User, *store.Error) {

	tx, err := repo.store.db.Begin()
	if err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	query := "UPDATE users SET balance = $2 WHERE user_id = $1 RETURNING user_id, balance"

	if err := tx.QueryRow(query, u.Id, balance).Scan(&u.Id, &u.Balance); err != nil {
		tx.Rollback()
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	tx.Commit()

	return u, nil
}
