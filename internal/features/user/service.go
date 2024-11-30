package user

import (
	"context"
	"go02/internal/package/apperrors"
	"go02/internal/package/db"
	"go02/internal/package/logging"
	"log/slog"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type UserService interface {
	CreateUser(ctx context.Context, name string, age int, bio string, avatarURL string) error
	UpdateUser(ctx context.Context, ID int, name string, age int, bio string, avatarURL string) error
	DeleteUser(ctx context.Context, ID int) error
	GetUserList(ctx context.Context, limit int, offset int) (GetUserListResponse, error)
	GetUserOne(ctx context.Context, ID int) (GetUserResponse, error)
}

type userService struct {
	transaction    db.Transaction
	userRepository UserRepository
}

func NewUserService(
	transaction db.Transaction,
	userRepository UserRepository,
) UserService {
	return &userService{
		transaction:    transaction,
		userRepository: userRepository,
	}
}

func (u *userService) CreateUser(ctx context.Context, name string, age int, bio string, avatarURL string) error {

	err := u.transaction.WithinTransaction(ctx, func(ctx context.Context) error {
		user, err := NewUser(name, age)
		if err != nil {
			return apperrors.WithStack(err)
		}

		_, err = u.userRepository.Create(ctx, user)
		if err != nil {
			return apperrors.WithStack(err)
		}

		profile, err := NewProfile(user.ID, bio, avatarURL)
		if err != nil {
			return apperrors.WithStack(err)
		}

		_, err = u.userRepository.CreateProfile(ctx, profile)
		if err != nil {
			return apperrors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (u *userService) UpdateUser(ctx context.Context, ID int, name string, age int, bio string, avatarURL string) error {

	err := u.transaction.WithinTransaction(ctx, func(ctx context.Context) error {
		user, err := u.userRepository.GetOne(ctx, ID)
		if err != nil {
			return apperrors.WithStack(err)
		}

		user.Name = name
		user.Age = age

		err = u.userRepository.Update(ctx, &user)
		if err != nil {
			return apperrors.WithStack(err)
		}

		profile, err := u.userRepository.GetProfileByUserID(ctx, ID)
		if err != nil {
			return apperrors.WithStack(err)
		}

		profile.Bio = bio
		profile.AvatarURL = avatarURL

		err = u.userRepository.UpdateProfile(ctx, &profile)
		if err != nil {
			return apperrors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (u *userService) DeleteUser(ctx context.Context, ID int) error {

	err := u.userRepository.Delete(ctx, ID)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (u *userService) GetUserList(ctx context.Context, limit int, offset int) (GetUserListResponse, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "userService.GetUserList")
	defer span.End()

	resUsers := GetUserListResponse{
		Users: []GetUserResponse{},
	}

	// limit default is 100
	l := limit
	if l == 0 {
		l = 100
	}

	users, err := u.userRepository.GetList(ctx, l, offset)
	if err != nil {
		return resUsers, apperrors.WithStack(err)
	}

	resUsers = GetUserListResponse{
		Users: lo.Map(users, func(u User, _ int) GetUserResponse {
			return GetUserResponse{
				ID:   u.ID,
				Name: u.Name,
				Age:  u.Age,
			}
		}),
	}

	logging.Info(ctx, "success to get user list", slog.Any("users", resUsers))

	return resUsers, nil
}

func (u *userService) GetUserOne(ctx context.Context, ID int) (GetUserResponse, error) {
	var resUser GetUserResponse

	user, err := u.userRepository.GetOne(ctx, ID)
	if err != nil {
		return resUser, apperrors.WithStack(err)
	}

	resUser = GetUserResponse{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
	}

	return resUser, nil
}
