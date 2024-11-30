package user

import (
	"context"
	"go02/internal/package/apperrors"
	"go02/internal/package/db"

	"github.com/cockroachdb/errors"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Repository interface {
	CreateUser(ctx context.Context, data *User) (int, error)
	UpdateUser(ctx context.Context, data *User) error
	DeleteUser(ctx context.Context, userID int) error
	GetUserList(ctx context.Context, limit int, offset int) ([]User, error)
	GetUserOne(ctx context.Context, userID int) (User, error)

	CreateProfile(ctx context.Context, data *Profile) (int, error)
	UpdateProfile(ctx context.Context, data *Profile) error
	DeleteProfile(ctx context.Context, userID int) error
	GetProfileByUserID(ctx context.Context, userID int) (Profile, error)
}

type repository struct {
	conn *bun.DB
}

func NewRepository(conn *bun.DB) Repository {
	return &repository{
		conn: conn,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (int, error) {

	tx := db.GetTxOrDB(ctx, r.conn)
	_, err := tx.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return 0, apperrors.WithStack(err)
	}

	return user.ID, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) error {

	tx := db.GetTxOrDB(ctx, r.conn)
	_, err := tx.NewUpdate().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (r *repository) DeleteUser(ctx context.Context, userID int) error {
	_, err := r.conn.NewDelete().Model(&User{}).Where("id = ?", userID).Exec(ctx)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (r *repository) GetUserList(ctx context.Context, limit int, offset int) ([]User, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "userRepository.GetUserList")
	defer span.End()

	span.SetAttributes(attribute.String("db.operation", "select"))
	span.SetAttributes(attribute.String("db.table", "users"))

	users := make([]User, 0, limit)

	if err := r.conn.NewSelect().Model(&users).Limit(limit).Offset(offset).Scan(ctx); err != nil {
		return []User{}, errors.WithStack(err)
	}

	if len(users) == 0 {
		return []User{}, apperrors.WithStack(apperrors.ErrNotFound)
	}

	return users, nil
}

func (r *repository) GetUserOne(ctx context.Context, userID int) (User, error) {
	var user User

	if err := r.conn.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx); err != nil {
		return User{}, apperrors.WithStack(err)
	}

	return user, nil
}

func (r *repository) CreateProfile(ctx context.Context, profile *Profile) (int, error) {

	tx := db.GetTxOrDB(ctx, r.conn)
	_, err := tx.NewInsert().Model(profile).Exec(ctx)
	if err != nil {
		return 0, apperrors.WithStack(err)
	}

	return profile.ID, nil
}

func (r *repository) UpdateProfile(ctx context.Context, profile *Profile) error {

	tx := db.GetTxOrDB(ctx, r.conn)
	_, err := tx.NewUpdate().Model(profile).WherePK().Exec(ctx)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (r *repository) DeleteProfile(ctx context.Context, profileID int) error {
	_, err := r.conn.NewDelete().Model(&Profile{}).Where("id = ?", profileID).Exec(ctx)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (r *repository) GetProfileByUserID(ctx context.Context, userID int) (Profile, error) {
	var profile Profile

	if err := r.conn.NewSelect().Model(&profile).Where("user_id = ?", userID).Scan(ctx); err != nil {
		return Profile{}, apperrors.WithStack(err)
	}

	return profile, nil
}
