package repository

import (
	"context"
	"go02/model"

	"github.com/uptrace/bun"
)

type UserRepository interface {
	Create(ctx context.Context, data *model.User) (int, error)
	Update(ctx context.Context, data *model.User) error
	Delete(ctx context.Context, userID int) error
	GetList(ctx context.Context, limit int, offset int) ([]model.User, error)
	GetOne(ctx context.Context, userID int) (model.User, error)
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create Userの新規作成
func (r *userRepository) Create(ctx context.Context, user *model.User) (int, error) {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Update Userの更新
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.NewUpdate().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Delete Userの削除
func (r *userRepository) Delete(ctx context.Context, userID int) error {
	_, err := r.db.NewDelete().Model(&model.User{}).Where("id = ?", userID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetList Userの複数件取得
func (r *userRepository) GetList(ctx context.Context, limit int, offset int) ([]model.User, error) {
	users := make([]model.User, 0, limit)

	if err := r.db.NewSelect().Model(&users).Limit(limit).Offset(offset).Scan(ctx); err != nil {
		return users, err
	}

	return users, nil
}

// GetOne Userを1件取得
func (r *userRepository) GetOne(ctx context.Context, userID int) (model.User, error) {
	var user model.User

	if err := r.db.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx); err != nil {
		return user, err
	}

	return user, nil
}
