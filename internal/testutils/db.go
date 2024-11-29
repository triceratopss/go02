package testutils

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type pgtestcontainers struct {
	postgres.PostgresContainer
	DSN      string
	TearDown func()
}

func PrepareContainer(ctx context.Context, t *testing.T) pgtestcontainers {
	t.Helper()

	dbName := "your-db-name"
	dbUser := "root"
	dbPassword := "root"

	container, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	teardown := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	return pgtestcontainers{
		PostgresContainer: *container,
		DSN:               dsn,
		TearDown:          teardown,
	}
}

func OpenDBForTest(t *testing.T, dsn string) (*bun.DB, error) {
	t.Helper()

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	db := bun.NewDB(sqlDB, pgdialect.New())

	return db, nil
}

func MigrateUp(t *testing.T, dsn string) error {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	migrationsPath := filepath.Join(projectRoot, "db-migration", "migrations")
	sourceURL := "file://" + migrationsPath

	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
