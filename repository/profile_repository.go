package repository

import (
	"context"
	"go02/model"
	"go02/packages/apperrors"
	"go02/packages/db"

	"github.com/uptrace/bun"
)

type ProfileRepository interface {
	Create(ctx context.Context, data *model.Profile) (int, error)
	Update(ctx context.Context, data *model.Profile) error
	Delete(ctx context.Context, userID int) error
	GetProfileByUserID(ctx context.Context, userID int) (model.Profile, error)
}

type profileRepository struct {
	conn *bun.DB
}

func NewProfileRepository(conn *bun.DB) ProfileRepository {
	return &profileRepository{
		conn: conn,
	}
}

func (r *profileRepository) Create(ctx context.Context, profile *model.Profile) (int, error) {

	tx := db.GetTxOrDB(ctx, r.conn)
	_, err := tx.NewInsert().Model(profile).Exec(ctx)
	if err != nil {
		return 0, apperrors.WithStack(err)
	}

	return profile.ID, nil
}

func (r *profileRepository) Update(ctx context.Context, profile *model.Profile) error {

	tx := db.GetTxOrDB(ctx, r.conn)
	_, err := tx.NewUpdate().Model(profile).WherePK().Exec(ctx)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (r *profileRepository) Delete(ctx context.Context, profileID int) error {
	_, err := r.conn.NewDelete().Model(&model.Profile{}).Where("id = ?", profileID).Exec(ctx)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (r *profileRepository) GetProfileByUserID(ctx context.Context, userID int) (model.Profile, error) {
	var profile model.Profile

	if err := r.conn.NewSelect().Model(&profile).Where("user_id = ?", userID).Scan(ctx); err != nil {
		return model.Profile{}, apperrors.WithStack(err)
	}

	return profile, nil
}
