package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"go02/usecase"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	CreateUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	GetUserList(c echo.Context) error
	GetUserOne(c echo.Context) error
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var params struct {
		Name      string `json:"name"`
		Age       int    `json:"age"`
		Bio       string `json:"bio"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "bad request",
		})
	}

	err := h.userUsecase.CreateUser(ctx, params.Name, params.Age, params.Bio, params.AvatarURL)
	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "failed to create user",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "invalid id",
		})
	}

	var params struct {
		Name      string `json:"name"`
		Age       int    `json:"age"`
		Bio       string `json:"bio"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "bad request",
		})
	}

	if err := h.userUsecase.UpdateUser(ctx, id, params.Name, params.Age, params.Bio, params.AvatarURL); err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "failed to update user",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "invalid id",
		})
	}

	if err := h.userUsecase.DeleteUser(ctx, id); err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "failed to delete user",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "success",
	})
}

func (h *userHandler) GetUserList(c echo.Context) error {
	ctx := c.Request().Context()

	var params struct {
		Limit  int `query:"name"`
		Offset int `query:"age"`
	}

	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "bad request",
		})
	}

	resUsers, err := h.userUsecase.GetUserList(ctx, params.Limit, params.Offset)
	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, resUsers)
}

func (h *userHandler) GetUserOne(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "invalid id",
		})
	}

	resUser, err := h.userUsecase.GetUserOne(ctx, id)
	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "failed to get user",
		})
	}

	return c.JSON(http.StatusOK, resUser)
}
