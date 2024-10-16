package usecase

import (
	"context"
	"go02/model"
	"go02/packages/apperrors"
	"go02/repository"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

// UserUsecase User 関係のusecaseのinterface
type UserUsecase interface {
	CreateUser(ctx context.Context, name string, age int, bio string, avatarURL string) error
	UpdateUser(ctx context.Context, ID int, name string, age int, bio string, avatarURL string) error
	DeleteUser(ctx context.Context, ID int) error
	GetUserList(ctx context.Context, limit int, offset int) (ResGetUserList, error)
	GetUserOne(ctx context.Context, ID int) (ResGetUser, error)
}

type userUsecase struct {
	transactionRepository repository.TransactionRepository
	userRepository        repository.UserRepository
	profileRepository     repository.ProfileRepository
}

// NewUserUsecase User usecaseのコンストラクタ
func NewUserUsecase(
	transactionRepository repository.TransactionRepository,
	userRepository repository.UserRepository,
	profileRepository repository.ProfileRepository,
) UserUsecase {
	return &userUsecase{
		transactionRepository: transactionRepository,
		userRepository:        userRepository,
		profileRepository:     profileRepository,
	}
}

type ReqCreateUpdateUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type ReqGetUserList struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
type ResGetUserList struct {
	Users []ResGetUser `json:"users"`
}
type ResGetUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *userUsecase) CreateUser(ctx context.Context, name string, age int, bio string, avatarURL string) error {

	err := u.transactionRepository.WithinTransaction(ctx, func(ctx context.Context) error {
		user, err := model.NewUser(name, age)
		if err != nil {
			return apperrors.WithStack(err)
		}

		_, err = u.userRepository.Create(ctx, user)
		if err != nil {
			return apperrors.WithStack(err)
		}

		profile, err := model.NewProfile(user.ID, bio, avatarURL)
		if err != nil {
			return apperrors.WithStack(err)
		}

		_, err = u.profileRepository.Create(ctx, profile)
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

func (u *userUsecase) UpdateUser(ctx context.Context, ID int, name string, age int, bio string, avatarURL string) error {

	err := u.transactionRepository.WithinTransaction(ctx, func(ctx context.Context) error {
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

		profile, err := u.profileRepository.GetProfileByUserID(ctx, ID)
		if err != nil {
			return apperrors.WithStack(err)
		}

		profile.Bio = bio
		profile.AvatarURL = avatarURL

		err = u.profileRepository.Update(ctx, &profile)
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

func (u *userUsecase) DeleteUser(ctx context.Context, ID int) error {

	err := u.userRepository.Delete(ctx, ID)
	if err != nil {
		return apperrors.WithStack(err)
	}

	return nil
}

func (u *userUsecase) GetUserList(ctx context.Context, limit int, offset int) (ResGetUserList, error) {
	tracer := otel.Tracer("usecase")
	ctx, span := tracer.Start(ctx, "userUsecase.GetUserList")
	defer span.End()

	resUsers := ResGetUserList{
		Users: []ResGetUser{},
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

	resUsers = ResGetUserList{
		Users: lo.Map(users, func(u model.User, _ int) ResGetUser {
			return ResGetUser{
				ID:   u.ID,
				Name: u.Name,
				Age:  u.Age,
			}
		}),
	}

	return resUsers, nil
}

func (u *userUsecase) GetUserOne(ctx context.Context, ID int) (ResGetUser, error) {
	var resUser ResGetUser

	user, err := u.userRepository.GetOne(ctx, ID)
	if err != nil {
		return resUser, apperrors.WithStack(err)
	}

	resUser = ResGetUser{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
	}

	return resUser, nil
}
