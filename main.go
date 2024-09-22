package main

import (
	"go02/interface/router"
	"go02/middleware"
	"go02/packages/config"
	"go02/packages/db"
	"go02/packages/logging"
	"log"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}

func run() error {
	err := config.Init()
	if err != nil {
		return errors.Wrap(err, "failed to initialize config")
	}

	logging.Init()

	db, err := db.OpenDB()
	if err != nil {
		return errors.Wrap(err, "failed to initialize a new database")
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())

	router.Init(e, db)

	port := "8080"
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: e,
	}

	logging.Infof("listening on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		return errors.Wrap(err, "failed to listen and serve")
	}

	return nil
}
