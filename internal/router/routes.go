package router

import (
	"go02/internal/handler"
	"go02/internal/repository"
	"go02/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func Init(e *echo.Echo, db *bun.DB) {

	transactionRepository := repository.NewTransactionRepository(db)
	userRepository := repository.NewUserRepository(db)
	profileRepository := repository.NewProfileRepository(db)
	userService := service.NewUserService(transactionRepository, userRepository, profileRepository)
	userHandler := handler.NewUserHandler(userService)

	e.POST("/users", userHandler.CreateUser)
	e.GET("/users", userHandler.GetUserList)
	e.GET("/users/:id", userHandler.GetUserOne)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)
}
