package db

import (
	"context"
	"database/sql"
	"fmt"
	"go02/packages/config"
	"net/url"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	_ "github.com/lib/pq"
)

func OpenDB() (*bun.DB, error) {

	encodedDBUser := url.QueryEscape(config.Config.DBUser)
	encodedDBPassword := url.QueryEscape(config.Config.DBPassword)
	encodedDBName := url.QueryEscape(config.Config.DBName)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		encodedDBUser,
		encodedDBPassword,
		config.Config.DBHost,
		config.Config.DBPort,
		encodedDBName,
	)
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

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
