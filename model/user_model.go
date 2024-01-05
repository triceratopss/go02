package model

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
