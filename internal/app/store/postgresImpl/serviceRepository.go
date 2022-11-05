package sqlstore

import (
	"balance-microservice/internal/app/model"
	"balance-microservice/internal/app/store"
	"database/sql"
	"fmt"
)

type ServiceRepository struct {
	store *Store
}

func (repo *ServiceRepository) FindById(id int) (*model.Service, *store.Error) {
	query := "SELECT * FROM services WHERE service_id = $1"

	s := &model.Service{}

	if err := repo.store.db.QueryRow(
		query,
		id,
	).Scan(&s.Id, &s.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, &store.Error{
				Message:    fmt.Sprintf("Service with id %d is not found", id),
				StatusCode: 404,
			}
		}

		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return s, nil
}
