package model

import (
	"time"

	"github.com/uptrace/bun"
)

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
