package usecase

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type GetAccountOutput struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
	AvatarURL   *string
	CreateTime  time.Time
}
