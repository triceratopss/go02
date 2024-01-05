package usecase

import (
	"context"
	"go02/model"
	"go02/repository"

	"github.com/samber/lo"
)

// UserUsecase User 関係のusecaseのinterface
type UserUsecase interface {
	CreateUser(ctx context.Context, name string, age int) error
	UpdateUser(ctx context.Context, ID int, name string, age int) error
	DeleteUser(ctx context.Context, ID int) error
	GetUserList(ctx context.Context, limit int, offset int) (ResGetUserList, error)
	GetUserOne(ctx context.Context, ID int) (ResGetUser, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase User usecaseのコンストラクタ
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
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

// CreateUser User を追加する
func (u *userUsecase) CreateUser(ctx context.Context, name string, age int) error {

	user, err := model.NewUser(name, age)
	if err != nil {
		return err
	}

	_, err = u.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser User を更新する
func (u *userUsecase) UpdateUser(ctx context.Context, ID int, name string, age int) error {

	user := model.User{
		ID:   ID,
		Name: name,
		Age:  age,
	}

	err := u.userRepo.Update(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

// // DeleteUser User を削除する
func (u *userUsecase) DeleteUser(ctx context.Context, ID int) error {

	err := u.userRepo.Delete(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}

// GetUserList User を複数件取得する
func (u *userUsecase) GetUserList(ctx context.Context, limit int, offset int) (ResGetUserList, error) {
	var resUsers ResGetUserList

	// limit default is 100
	l := limit
	if l == 0 {
		l = 100
	}

	users, err := u.userRepo.GetList(ctx, l, offset)
	if err != nil {
		return resUsers, err
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

// // GetUserList User を1件取得する
func (u *userUsecase) GetUserOne(ctx context.Context, ID int) (ResGetUser, error) {
	var resUser ResGetUser

	user, err := u.userRepo.GetOne(ctx, ID)
	if err != nil {
		return resUser, err
	}

	resUser = ResGetUser{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
	}

	return resUser, nil
}
