package handler

import (
	"errors"
	"net/http"
	"strconv"

	"go02/packages/apperrors"
	"go02/packages/logging"
	"go02/usecase"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
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
		logging.Errorf(ctx, err, "failed to bind request body: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "bad request",
		})
	}

	err := h.userUsecase.CreateUser(ctx, params.Name, params.Age, params.Bio, params.AvatarURL)
	if err != nil {
		logging.Errorf(ctx, err, "failed to CreateUser: %s", err.Error())
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
		logging.Errorf(ctx, err, "failed to parse id: %s", err.Error())
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
		logging.Errorf(ctx, err, "failed to bind request body: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "bad request",
		})
	}

	if err := h.userUsecase.UpdateUser(ctx, id, params.Name, params.Age, params.Bio, params.AvatarURL); err != nil {
		logging.Errorf(ctx, err, "failed to UpdateUser: %s", err.Error())
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
		logging.Errorf(ctx, err, "failed to parse id: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "invalid id",
		})
	}

	if err := h.userUsecase.DeleteUser(ctx, id); err != nil {
		logging.Errorf(ctx, err, "failed to DeleteUser: %s", err.Error())
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
	tracer := otel.Tracer("handler")
	ctx, span := tracer.Start(ctx, "GetUserList")
	defer span.End()

	var params struct {
		Limit  int `query:"limit"`
		Offset int `query:"offset"`
	}

	if err := c.Bind(&params); err != nil {
		logging.Errorf(ctx, err, "failed to bind query params: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": http.StatusText(http.StatusBadRequest),
		})
	}

	resUsers, err := h.userUsecase.GetUserList(ctx, params.Limit, params.Offset)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		logging.Errorf(ctx, err, "failed to GetUserList: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"message": http.StatusText(http.StatusInternalServerError),
		})
	}

	return c.JSON(http.StatusOK, resUsers)
}

func (h *userHandler) GetUserOne(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logging.Errorf(ctx, err, "failed to parse id: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "invalid id",
		})
	}

	resUser, err := h.userUsecase.GetUserOne(ctx, id)
	if err != nil {
		logging.Errorf(ctx, err, "failed to GetUserOne: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "failed to get user",
		})
	}

	return c.JSON(http.StatusOK, resUser)
}
