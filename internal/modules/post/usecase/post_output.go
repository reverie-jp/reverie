package usecase

import (
	"time"

	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type PostOutput struct {
	ID            ulid.ULID
	ShortID       string
	Text          string
	Author        *usergw.UserView
	ReplyToID     *ulid.ULID
	RepostID      *ulid.ULID
	ReplyCount    int64
	RepostCount   int64
	FavoriteCount int64
	IsFavorited   bool
	CreateTime    time.Time
	RepostOf      *PostOutput
}
