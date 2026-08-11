package transaction

import "gorm.io/gorm"

type ITransactionManager interface {
	Execute(fn func(tx *gorm.DB) error) error
}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) ITransactionManager {
	if db == nil {
		panic("database connection (gorm.DB) cannot be nil in TransactionManager")
	}
	return &transactionManager{db: db}
}

func (t *transactionManager) Execute(fn func(tx *gorm.DB) error) error {
	return t.db.Transaction(fn)
}
