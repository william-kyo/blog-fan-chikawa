package repository

import (
	"blog-fanchiikawa-service/db"
)

// transactionManager implements TransactionManager interface
type transactionManager struct{}

// NewTransactionManager creates a new TransactionManager instance
func NewTransactionManager() TransactionManager {
	return &transactionManager{}
}

// WithTransaction executes a function within a database transaction
func (tm *transactionManager) WithTransaction(fn func() error) error {
	session := db.Engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	if err := fn(); err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}