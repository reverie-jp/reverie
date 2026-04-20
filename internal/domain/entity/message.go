package entity

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type Message struct {
	ID         ulid.ULID
	RoomID     ulid.ULID
	SenderID   ulid.ULID
	Content    *string
	IsDeleted  bool
	IsEdited   bool
	CreateTime time.Time
}
