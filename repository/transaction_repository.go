package repository

import (
	"context"
	"database/sql"
	"go02/packages/apperrors"
	"go02/packages/db"

	"github.com/cockroachdb/errors"
	"github.com/uptrace/bun"
)

type TransactionRepository interface {
	WithinTransaction(ctx context.Context, f func(ctx context.Context) error) (err error)
}

func NewTransactionRepository(conn *bun.DB) TransactionRepository {
	return &transactionRepository{
		conn: conn,
	}
}

type transactionRepository struct {
	conn *bun.DB
}

func (r *transactionRepository) WithinTransaction(ctx context.Context, f func(ctx context.Context) error) (err error) {
	tx, err := db.GetTxOrDB(ctx, r.conn).BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return apperrors.WithStack(err)
	}

	defer func() {
		if r := recover(); r != nil {
			err = apperrors.WithStack(errors.Errorf("panic error %v", r))
			tx.Rollback()
		}
	}()

	err = f(db.SetTx(ctx, &tx))
	if err != nil {
		tx.Rollback()
		return apperrors.WithStack(err)
	}

	return apperrors.WithStack(tx.Commit())
}
