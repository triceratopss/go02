package repository

import (
	"context"
	"go02/internal/model"
	"go02/internal/package/apperrors"
	"go02/internal/package/db"

	"github.com/cockroachdb/errors"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserRepository interface {
	Create(ctx context.Context, data *model.User) (int, error)
	Update(ctx context.Context, data *model.User) error
	Delete(ctx context.Context, userID int) error
	GetList(ctx context.Context, limit int, offset int) ([]model.User, error)
	GetOne(ctx context.Context, userID int) (model.User, error)
}

type userRepository struct {
	conn *bun.DB
}

func NewUserRepository(conn *bun.DB) UserRepository {
	return &userRepository{
		conn: conn,
	}
}

// Create Userの新規作成
func (r *userRepository) Create(ctx context.Context, user *model.User) (int, error) {

	tx := db.GetTxOrDB(ctx, r.conn)
	_, err := tx.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return 0, apperrors.WithStack(err)
	}

	return user.ID, nil
}

// Update Userの更新
func (r *userRepository) Update(ctx context.Context, user *model.User) error {

	tx := db.GetTxOrDB(ctx, r.conn)
	_, err := tx.NewUpdate().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

// Delete Userの削除
func (r *userRepository) Delete(ctx context.Context, userID int) error {
	_, err := r.conn.NewDelete().Model(&model.User{}).Where("id = ?", userID).Exec(ctx)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

// GetList Userの複数件取得
func (r *userRepository) GetList(ctx context.Context, limit int, offset int) ([]model.User, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "userRepository.GetList")
	defer span.End()

	span.SetAttributes(attribute.String("db.operation", "select"))
	span.SetAttributes(attribute.String("db.table", "users"))

	users := make([]model.User, 0, limit)

	if err := r.conn.NewSelect().Model(&users).Limit(limit).Offset(offset).Scan(ctx); err != nil {
		return []model.User{}, errors.WithStack(err)
	}

	if len(users) == 0 {
		return []model.User{}, apperrors.WithStack(apperrors.ErrNotFound)
	}

	return users, nil
}

// GetOne Userを1件取得
func (r *userRepository) GetOne(ctx context.Context, userID int) (model.User, error) {
	var user model.User

	if err := r.conn.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx); err != nil {
		return model.User{}, apperrors.WithStack(err)
	}

	return user, nil
}
