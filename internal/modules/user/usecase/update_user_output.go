package usecase

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type UpdateUserOutput struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
	Biography   string
	IsPrivate   bool
	CreateTime  time.Time
	UpdateTime  time.Time
}
