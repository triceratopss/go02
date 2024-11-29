package main

import (
	"context"
	"go02/interface/router"
	"go02/middleware"
	"go02/packages/config"
	"go02/packages/db"
	"go02/packages/logging"
	"go02/packages/tracer"
	"log"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}

func run() error {
	ctx := context.Background()

	err := config.Init()
	if err != nil {
		return errors.Wrap(err, "failed to initialize config")
	}

	logging.Init()

	db, err := db.OpenDB()
	if err != nil {
		return errors.Wrap(err, "failed to initialize a new database")
	}

	tp := tracer.InitializeTracer()
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}()

	e := echo.New()

	e.Use(otelecho.Middleware("go02"))
	e.Use(middleware.Logger())

	router.Init(e, db)

	port := "8080"
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: e,
	}

	logging.Infof(ctx, "listening on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		return errors.Wrap(err, "failed to listen and serve")
	}

	return nil
}
