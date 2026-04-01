package entity

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type User struct {
	ID                ulid.ULID
	CustomID          string
	CustomIDChangedAt *time.Time
	DisplayName       string
	Biography         *string
	Location          *string
	Website           *string
	AvatarURL         *string
	BannerURL         *string
	IsPrivate         bool
	CreateTime        time.Time
	UpdateTime        time.Time
}

type AuthProvider struct {
	UserID ulid.ULID
}
