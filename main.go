package main

import (
	"fmt"
	"go02/interface/router"
	"go02/middleware"
	"go02/packages/db"
	"go02/packages/logger"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	logger.Init()

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

	slog.Info(fmt.Sprintf("listening on port %s", port))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
