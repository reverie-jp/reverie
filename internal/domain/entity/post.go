package entity

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type Post struct {
	ID         ulid.ULID
	AuthorID   ulid.ULID
	ShortID    string
	ReplyToPostID  *ulid.ULID
	RepostPostID   *ulid.ULID
	Text       string
	CreateTime time.Time
	UpdateTime time.Time
}
