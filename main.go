package main

import (
	"go02/interface/router"
	"go02/packages/db"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

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

	log.Printf("listening on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
