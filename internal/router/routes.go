package router

import (
	"go02/internal/features/user"
	"go02/internal/repository"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func Init(e *echo.Echo, db *bun.DB) {

	transactionRepository := repository.NewTransactionRepository(db)
	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(transactionRepository, userRepository)
	userHandler := user.NewUserHandler(userService)

	e.POST("/users", userHandler.CreateUser)
	e.GET("/users", userHandler.GetUserList)
	e.GET("/users/:id", userHandler.GetUserOne)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)
}
