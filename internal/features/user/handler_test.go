package user_test

import (
	"context"
	"go02/internal/features/user"
	"go02/internal/repository"
	"go02/internal/testutils"
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
		wantError        bool
		expectedStatus   int
		expectedFilePath string
		testData         []user.User
	}{
		{
			name:             "正常系: アイテムが存在する場合",
			queryParams:      nil,
			wantError:        false,
			expectedStatus:   http.StatusOK,
			expectedFilePath: "testdata/get_users/ok_res.golden.json",
			testData: []user.User{
				{ID: 1, Name: "taro", Age: 24},
				{ID: 2, Name: "takeshi", Age: 20},
			},
		},
		{
			name:             "正常系: アイテムが存在しない場合",
			queryParams:      nil,
			wantError:        false,
			expectedStatus:   http.StatusOK,
			expectedFilePath: "testdata/get_users/ok_res_empty.golden.json",
			testData:         []user.User{},
		},
		{
			name:             "正常系: アイテムが存在し、クエリパラメータが指定されている場合",
			queryParams:      map[string]string{"limit": "3", "offset": "2"},
			wantError:        false,
			expectedStatus:   http.StatusOK,
			expectedFilePath: "testdata/get_users/ok_res_query_param.golden.json",
			testData: []user.User{
				{ID: 1, Name: "taro", Age: 24},
				{ID: 2, Name: "takeshi", Age: 20},
				{ID: 3, Name: "hanako", Age: 21},
				{ID: 4, Name: "kana", Age: 27},
				{ID: 5, Name: "yuki", Age: 18},
				{ID: 6, Name: "ichiro", Age: 30},
			},
		},
		{
			name:             "異常系: クエリパラメータの値が不正な場合",
			queryParams:      map[string]string{"limit": "a", "offset": "b"},
			wantError:        true,
			expectedStatus:   http.StatusBadRequest,
			expectedFilePath: "testdata/get_users/err_res_400.golden.json",
			testData:         []user.User{},
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
			req := httptest.NewRequest(http.MethodGet, "/users?"+q.Encode(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			transactionRepository := repository.NewTransactionRepository(db)
			userRepository := user.NewUserRepository(db)
			userService := user.NewUserService(transactionRepository, userRepository)
			userHandler := user.NewUserHandler(userService)

			// Act
			err = userHandler.GetUserList(c)

			// Assert
			var actualJSON string

			if tt.wantError {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)

				errorResponse, ok := he.Message.(map[string]any)
				assert.True(t, ok)
				actualJSON, err = testutils.MapToJSONString(errorResponse)
				assert.NoError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)

				actualJSON = rec.Body.String()
			}

			expectedJSON, err := testutils.ReadJSONFile(t, tt.expectedFilePath)
			if err != nil {
				t.Fatalf("Failed to read expected JSON file: %v", err)
			}

			assert.JSONEq(t, expectedJSON, actualJSON)
		})
	}
}
