package sqlstore

import (
	"balance-microservice/internal/app/model"
	"balance-microservice/internal/app/store"
	"balance-microservice/internal/app/utils/csv"
	"fmt"
	"time"
)

type AccountingRepository struct {
	store *Store
}

func (repo *AccountingRepository) AddRevenue(req *model.RevenueAddRequest) *store.Error {

	tx, err := repo.store.db.Begin()
	if err != nil {
		return &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	user, e := repo.store.UserRepository().FindById(req.UserId)
	if e != nil {
		return e
	}

	fmt.Printf("User: %v", user)

	service, e := repo.store.ServiceRepository().FindById(req.ServiceId)
	if e != nil {
		return e
	}

	accounting := &model.Accounting{
		UserId:  user.Id,
		Service: service,
		OrderId: req.OrderId,
		Price:   req.Price,
	}

	query := "INSERT INTO accounting (user_id, service_id, order_id, price, created_at) " +
		"VALUES($1, $2, $3, $4, now()) RETURNING record_id, created_at"

	if err := tx.QueryRow(query, req.UserId, service.Id, req.OrderId, req.Price).
		Scan(&accounting.Id, &accounting.CreatedAt); err != nil {
		tx.Rollback()
		return &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	tx.Commit()

	return nil
}

func (repo *AccountingRepository) FindAll() ([]model.Accounting, *store.Error) {

	query := "SELECT * FROM accounting_record"

	records := make([]model.Accounting, 0)

	rows, err := repo.store.db.Query(query)
	if err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	for rows.Next() {

		var a model.Accounting
		var s model.Service

		if err := rows.Scan(&a.Id, &a.UserId, &a.OrderId, &a.Price, &a.CreatedAt, &s.Id, &s.Name); err != nil {
			return nil, &store.Error{
				Message:    err.Error(),
				StatusCode: 500,
			}
		}

		a.Service = &s

		records = append(records, a)
	}

	if err = rows.Err(); err != nil {
		return records, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	return records, nil
}

func (repo *AccountingRepository) ExportCSV(dt string) (*string, *store.Error) {
	filesPath := "./files/csv/"

	fmt.Println(dt)

	date, err := time.Parse("2006-01", dt)
	if err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 400,
		}
	}

	type csvStruct struct {
		ServiceName string
		Value       float64
	}

	records := make([][]string, 0)

	query := "SELECT s.name, sum(a.price) FROM accounting AS a " +
		"LEFT JOIN services as s on s.service_id = a.service_id " +
		"WHERE date_part('year', created_at) = $1 AND date_part('month', created_at) = $2 " +
		"GROUP BY s.name;"

	rows, err := repo.store.db.Query(query, date.Year(), date.Month())
	if err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	for rows.Next() {

		csvLine := &csvStruct{}

		if err := rows.Scan(&csvLine.ServiceName, &csvLine.Value); err != nil {
			return nil, &store.Error{
				Message:    err.Error(),
				StatusCode: 500,
			}
		}

		line := []string{csvLine.ServiceName, fmt.Sprintf("%.2f", csvLine.Value)}

		records = append(records, line)
	}

	if err = rows.Err(); err != nil {
		return nil, &store.Error{
			Message:    err.Error(),
			StatusCode: 500,
		}
	}

	filename := fmt.Sprintf("%d-%d", date.Year(), date.Month())

	if err := csv.WriteCVS(records, filename); err != nil {
		return nil, err
	}

	filePath := filesPath + filename + ".csv"

	return &filePath, nil
}
