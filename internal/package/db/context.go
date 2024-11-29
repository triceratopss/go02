package db

import (
	"context"

	"github.com/uptrace/bun"
)

type dbTx struct{}

func SetTx(ctx context.Context, tx *bun.Tx) context.Context {
	return context.WithValue(ctx, dbTx{}, tx)
}

func getTx(ctx context.Context) *bun.Tx {
	if tx, ok := ctx.Value(dbTx{}).(*bun.Tx); ok {
		return tx
	}
	return nil
}

func GetTxOrDB(ctx context.Context, db *bun.DB) bun.IDB {
	if tx := getTx(ctx); tx != nil {
		return tx
	}
	return db
}
