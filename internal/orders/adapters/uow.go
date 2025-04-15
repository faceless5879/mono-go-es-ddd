package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// UnitOfWork represents a database transaction that coordinates operations across repositories
type UnitOfWork struct {
	db         *sqlx.DB
	tx         *sqlx.Tx
	committed  bool
	rolledBack bool
}

// NewUnitOfWork creates a new UnitOfWork instance
func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
	return &UnitOfWork{
		db:         db,
		committed:  false,
		rolledBack: false,
	}
}

// Begin starts a new transaction, sqlx will trigger panic if error occurs
func (uow *UnitOfWork) Begin(ctx context.Context, options *sql.TxOptions) {
	tx := uow.db.MustBeginTx(ctx, options)
	uow.tx = tx
}

// Commit commits the transaction
func (uow *UnitOfWork) Commit() error {
	if uow.tx == nil {
		return fmt.Errorf("no transaction to commit")
	}

	err := uow.tx.Commit()
	if err == nil {
		uow.committed = true
	}
	return err
}

// Rollback rolls back the transaction
func (uow *UnitOfWork) Rollback() error {
	if uow.tx == nil {
		return fmt.Errorf("no transaction to rollback")
	}

	err := uow.tx.Rollback()
	if err == nil {
		uow.rolledBack = true
	}
	return err
}

// Execute will wrap execution in transaction and execution fn
func (uow *UnitOfWork) Execute(fn func(tx *sqlx.Tx) error) error {
	err := fn(uow.tx)
	if err != nil {
		rbErr := uow.Rollback()
		if rbErr != nil {
			return fmt.Errorf("execute error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return nil
}

// ExecuteWithResult will wrap execution in transaction and return T result
func ExecuteWithResult[T any](uow *UnitOfWork, fn func(tx *sqlx.Tx) (T, error)) (T, error) {
	var result T
	result, err := fn(uow.tx)
	if err != nil {
		rbErr := uow.Rollback()
		if rbErr != nil {
			return result, fmt.Errorf("execute error: %v, rollback error: %v", err, rbErr)
		}
		return result, err
	}
	return result, nil
}
