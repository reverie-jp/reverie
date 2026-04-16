package usecase

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type GetUserOutput struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
	Biography   string
	IsPrivate   bool
	IsMe        bool
	CreateTime  time.Time
}
