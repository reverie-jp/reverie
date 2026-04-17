package usecase

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type PostAuthorOutput struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
	IsPrivate   bool
}

type PostOutput struct {
	ID            ulid.ULID
	Text          string
	Author        *PostAuthorOutput
	ReplyToID     *ulid.ULID
	RepostID      *ulid.ULID
	ReplyCount    int64
	RepostCount   int64
	FavoriteCount int64
	IsFavorited   bool
	CreateTime    time.Time
	RepostOf      *PostOutput
}
