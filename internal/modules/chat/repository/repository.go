package repository

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type Repository interface {
	ListRoomsByUser(ctx context.Context, userID ulid.ULID) ([]*RoomRow, error)
	GetRoomByID(ctx context.Context, id ulid.ULID) (*entity.Room, error)
	GetDirectRoom(ctx context.Context, userID1, userID2 ulid.ULID) (ulid.ULID, error)
	CreateRoom(ctx context.Context, id ulid.ULID, roomType string, name *string) (*entity.Room, error)
	CreateRoomMember(ctx context.Context, id, roomID, userID ulid.ULID) error
	GetRoomMember(ctx context.Context, roomID, userID ulid.ULID) (*sqlc.RoomMember, error)
	GetRoomOtherMember(ctx context.Context, roomID, userID ulid.ULID) (ulid.ULID, error)
	GetMessageByID(ctx context.Context, id ulid.ULID) (*entity.Message, error)
	ListMessages(ctx context.Context, roomID ulid.ULID, cursor *time.Time, limit int32) ([]*entity.Message, error)
	CreateMessage(ctx context.Context, id, roomID, senderID ulid.ULID, content string) (*entity.Message, error)
	MarkRoomAsRead(ctx context.Context, roomID, userID ulid.ULID) error
	PinRoom(ctx context.Context, id, userID, roomID ulid.ULID) error
	UnpinRoom(ctx context.Context, userID, roomID ulid.ULID) error
	MuteRoom(ctx context.Context, roomID, userID ulid.ULID, muted bool) error
	LeaveRoom(ctx context.Context, roomID, userID ulid.ULID) error
	GetRoomMembers(ctx context.Context, roomID ulid.ULID) ([]ulid.ULID, error)
	RemoveRoomMember(ctx context.Context, roomID, userID ulid.ULID) error
	UpdateRoomName(ctx context.Context, roomID ulid.ULID, name *string) (*entity.Room, error)
	AddReaction(ctx context.Context, messageID, userID ulid.ULID, emoji string) error
	RemoveReaction(ctx context.Context, messageID, userID ulid.ULID, emoji string) error
	ListReactionsByMessage(ctx context.Context, messageID, userID ulid.ULID) ([]*ReactionRow, error)
	ListReactionsByMessages(ctx context.Context, messageIDs []ulid.ULID, userID ulid.ULID) ([]*ReactionByMessageRow, error)
}

type RoomRow struct {
	ID                  ulid.ULID
	RoomType            string
	Name                *string
	IsMuted             bool
	LastReadAt          *time.Time
	IsPinned            bool
	LastMessageText     *string
	LastMessageSenderID *string
	LastMessageAt       *time.Time
	UnreadCount         int64
}

type repositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &repositoryImpl{q: q}
}
