package handler_test

import (
	"context"
	"go02/interface/handler"
	"go02/model"
	"go02/repository"
	"go02/testutils"
	"go02/usecase"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestGetUserList(t *testing.T) {

	tests := []struct {
		name             string
		queryParams      map[string]string
		expectedStatus   int
		expectedFilePath string
		testData         []model.User
	}{
		{
			name:             "正常系: アイテムが存在する場合",
			queryParams:      nil,
			expectedStatus:   http.StatusOK,
			expectedFilePath: "testdata/get_users/ok_res.golden.json",
			testData: []model.User{
				{ID: 1, Name: "taro", Age: 24},
				{ID: 2, Name: "takeshi", Age: 20},
			},
		},
		{
			name:             "正常系: アイテムが存在しない場合",
			queryParams:      nil,
			expectedStatus:   http.StatusOK,
			expectedFilePath: "testdata/get_users/ok_res_empty.golden.json",
			testData:         []model.User{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			container := testutils.PrepareContainer(context.Background(), t)
			defer container.TearDown()

			db, err := testutils.OpenDBForTest(t, container.DSN)
			if err != nil {
				t.Fatal(err)
			}

			if err := testutils.MigrateUp(t, container.DSN); err != nil {
				t.Fatal(err)
			}

			testutils.PrepareTestDataForTestGetUserList(t, db, tt.testData)

			e := echo.New()

			q := make(url.Values)
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			req := httptest.NewRequest(http.MethodGet, "/users"+q.Encode(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			transactionRepository := repository.NewTransactionRepository(db)
			userRepository := repository.NewUserRepository(db)
			profileRepository := repository.NewProfileRepository(db)
			userUsecase := usecase.NewUserUsecase(transactionRepository, userRepository, profileRepository)
			userHandler := handler.NewUserHandler(userUsecase)

			// Act
			if err := userHandler.GetUserList(c); err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)

			actualJSON := rec.Body.String()

			expectedJSON, err := testutils.ReadJSONFile(t, tt.expectedFilePath)
			if err != nil {
				t.Fatalf("Failed to read expected JSON file: %v", err)
			}

			assert.JSONEq(t, expectedJSON, actualJSON)
		})
	}
}
