package router

import (
	"go02/internal/features/user"
	"go02/internal/package/db"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func Init(e *echo.Echo, bunDB *bun.DB) {

	transaction := db.NewTransaction(bunDB)
	userRepository := user.NewUserRepository(bunDB)
	userService := user.NewUserService(transaction, userRepository)
	userHandler := user.NewUserHandler(userService)

	e.POST("/users", userHandler.CreateUser)
	e.GET("/users", userHandler.GetUserList)
	e.GET("/users/:id", userHandler.GetUserOne)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)
}
