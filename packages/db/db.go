package db

import (
	"context"
	"database/sql"
	"fmt"
	"go02/packages/config"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	_ "github.com/lib/pq"
)

func OpenDB() (*bun.DB, error) {

	dsn := fmt.Sprintf("user=%s password=%s database=%s host=%s port=%s sslmode=disable",
		config.Config.DBUser,
		config.Config.DBPassword,
		config.Config.DBName,
		config.Config.DBHost,
		config.Config.DBPort,
	)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	db := bun.NewDB(sqlDB, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		// disable the hook
		bundebug.WithEnabled(false),

		// BUNDEBUG=1 logs failed queries
		// BUNDEBUG=2 logs all queries
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db, nil
}
