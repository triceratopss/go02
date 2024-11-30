package db

import (
	"context"
	"database/sql"
	"go02/internal/package/apperrors"

	"github.com/cockroachdb/errors"
	"github.com/uptrace/bun"
)

type Transaction interface {
	WithinTransaction(ctx context.Context, f func(ctx context.Context) error) (err error)
}

func NewTransaction(conn *bun.DB) Transaction {
	return &transaction{
		conn: conn,
	}
}

type transaction struct {
	conn *bun.DB
}

func (r *transaction) WithinTransaction(ctx context.Context, f func(ctx context.Context) error) (err error) {
	tx, err := GetTxOrDB(ctx, r.conn).BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return apperrors.WithStack(err)
	}

	defer func() {
		if r := recover(); r != nil {
			err = apperrors.WithStack(errors.Errorf("panic error %v", r))
			tx.Rollback()
		}
	}()

	err = f(SetTx(ctx, &tx))
	if err != nil {
		tx.Rollback()
		return apperrors.WithStack(err)
	}

	return apperrors.WithStack(tx.Commit())
}
