package main

import (
	"go02/interface/router"
	"go02/middleware"
	"go02/packages/db"
	"go02/packages/logging"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	logging.Init()

	db, err := db.OpenDB()
	if err != nil {
		log.Fatalf("failed to initialize a new database: %v", err)
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
		log.Fatal(err)
	}

}
