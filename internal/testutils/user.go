package testutils

import (
	"context"
	"go02/internal/features/user"
	"testing"

	"github.com/uptrace/bun"
)

func PrepareTestDataForTestGetUserList(t *testing.T, db *bun.DB, users []user.User) {
	t.Helper()

	if len(users) == 0 {
		return
	}

	ctx := context.Background()
	if _, err := db.NewInsert().Model(&users).Exec(ctx); err != nil {
		t.Fatal(err)
	}
}
