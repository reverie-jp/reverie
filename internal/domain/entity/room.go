package entity

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
)

type Room struct {
	ID         ulid.ULID
	RoomType   string
	Name       *string
	CreateTime time.Time
	UpdateTime time.Time
}

type RoomView struct {
	Room              *Room
	OtherUserID       *ulid.ULID
	LastMessageText   *string
	LastMessageSendID *ulid.ULID
	LastMessageAt     *time.Time
	UnreadCount       int64
	IsPinned          bool
	IsMuted           bool
}
