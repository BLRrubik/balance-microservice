package sqlstore

import (
	"balance-microservice/internal/app/store"
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	db                    *sql.DB
	userRepository        *UserRepository
	serviceRepository     *ServiceRepository
	billRepository        *BillRepository
	accountingRepository  *AccountingRepository
	transactionRepository *TransactionRepository
}

func (s *Store) TransactionRepository() store.TransactionRepository {
	if s.transactionRepository != nil {
		return s.transactionRepository
	}

	s.transactionRepository = &TransactionRepository{
		store: s,
	}

	return s.transactionRepository
}

func (s *Store) ServiceRepository() store.ServiceRepository {
	if s.serviceRepository != nil {
		return s.serviceRepository
	}

	s.serviceRepository = &ServiceRepository{
		store: s,
	}

	return s.serviceRepository
}

func (s *Store) AccountingRepository() store.AccountingRepository {
	if s.accountingRepository != nil {
		return s.accountingRepository
	}

	s.accountingRepository = &AccountingRepository{
		store: s,
	}

	return s.accountingRepository
}

func (s *Store) BillRepository() store.BillRepository {
	if s.billRepository != nil {
		return s.billRepository
	}

	s.billRepository = &BillRepository{
		store: s,
	}

	return s.billRepository
}

func (s *Store) UserRepository() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
