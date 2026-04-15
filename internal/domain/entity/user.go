package entity

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type User struct {
	ID                  ulid.ULID
	CustomID            string
	CustomIDChangedAt   *time.Time
	DisplayName         string
	Biography           *string
	AvatarMediaID       *ulid.ULID
	BannerMediaID       *ulid.ULID
	IsPrivate           bool
	Birthdate           *time.Time
	CreateTime          time.Time
	UpdateTime          time.Time
}

type AuthProvider struct {
	UserID ulid.ULID
}
