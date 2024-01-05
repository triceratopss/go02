package router

import (
	"go02/interface/handler"
	"go02/repository"
	"go02/usecase"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func Init(e *echo.Echo, db *bun.DB) {

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	e.POST("/users", userHandler.CreateUser)
	e.GET("/users", userHandler.GetUserList)
	e.GET("/users/:id", userHandler.GetUserOne)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)
}
