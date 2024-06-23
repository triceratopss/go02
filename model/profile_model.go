package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Profile struct {
	bun.BaseModel `bun:"table:profiles"`

	ID     int    `bun:",pk,autoincrement"`
	UserID int    `bun:"user_id"`
	Bio    string `bun:"bio"`

	CreatedAt time.Time `bun:",nullzero"`
	UpdatedAt time.Time `bun:",nullzero"`
}

func NewProfile(userID int, bio string) (*Profile, error) {

	profile := &Profile{
		UserID: userID,
		Bio:    bio,
	}

	return profile, nil
}
