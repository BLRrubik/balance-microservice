package store

type Store interface {
	UserRepository() UserRepository
	ServiceRepository() ServiceRepository
	BillRepository() BillRepository
	AccountingRepository() AccountingRepository
	TransactionRepository() TransactionRepository
}
