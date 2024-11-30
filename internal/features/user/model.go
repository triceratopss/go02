package user

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID   int    `bun:",pk,autoincrement"`
	Name string `bun:"name"`
	Age  int    `bun:"age"`

	CreatedAt time.Time `bun:",nullzero"`
	UpdatedAt time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
type Users []User

func NewUser(name string, age int) (*User, error) {

	user := &User{
		Name: name,
		Age:  age,
	}

	return user, nil
}

type Profile struct {
	bun.BaseModel `bun:"table:profiles"`

	ID        int    `bun:",pk,autoincrement"`
	UserID    int    `bun:"user_id"`
	Bio       string `bun:"bio"`
	AvatarURL string `bun:"avatar_url"`

	CreatedAt time.Time `bun:",nullzero"`
	UpdatedAt time.Time `bun:",nullzero"`
}

func NewProfile(userID int, bio string, avatarURL string) (*Profile, error) {

	profile := &Profile{
		UserID:    userID,
		Bio:       bio,
		AvatarURL: avatarURL,
	}

	return profile, nil
}
