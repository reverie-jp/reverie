package gateway

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/chat/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type RoomView struct {
	Room              *entity.Room
	OtherUser         *entity.User
	Members           []*entity.User
	LastMessageText   *string
	LastMessageSendID *ulid.ULID
	LastMessageAt     *time.Time
	UnreadCount       int64
	IsPinned          bool
	IsMuted           bool
}

type ReactionView struct {
	Emoji  string
	Count  int64
	IsMine bool
}

type MessageView struct {
	Message   *entity.Message
	Sender    *entity.User
	Reactions []*ReactionView
}

type Gateway interface {
	ListRooms(ctx context.Context, userID ulid.ULID) ([]*RoomView, error)
	GetOrCreateDirectRoom(ctx context.Context, myID, otherID ulid.ULID) (*RoomView, error)
	GetRoomView(ctx context.Context, roomID, userID ulid.ULID) (*RoomView, error)
	ListMessages(ctx context.Context, roomID ulid.ULID, cursor *time.Time, limit int32, requestorID ulid.ULID) ([]*MessageView, error)
	SendMessage(ctx context.Context, roomID, senderID ulid.ULID, content string) (*MessageView, error)
	MarkRoomAsRead(ctx context.Context, roomID, userID ulid.ULID) error
	PinRoom(ctx context.Context, userID, roomID ulid.ULID) (*RoomView, error)
	UnpinRoom(ctx context.Context, userID, roomID ulid.ULID) (*RoomView, error)
	MuteRoom(ctx context.Context, userID, roomID ulid.ULID) (*RoomView, error)
	UnmuteRoom(ctx context.Context, userID, roomID ulid.ULID) (*RoomView, error)
	LeaveRoom(ctx context.Context, roomID, userID ulid.ULID) error
	AddReaction(ctx context.Context, messageID, userID ulid.ULID, emoji string) (*MessageView, error)
	RemoveReaction(ctx context.Context, messageID, userID ulid.ULID, emoji string) (*MessageView, error)
	CreateGroupRoom(ctx context.Context, creatorID ulid.ULID, name string, memberIDs []ulid.ULID) (*RoomView, error)
	UpdateRoom(ctx context.Context, roomID ulid.ULID, name string) (*RoomView, error)
	ListRoomMembers(ctx context.Context, roomID ulid.ULID) ([]*entity.User, error)
	AddRoomMember(ctx context.Context, roomID, userID ulid.ULID) (*RoomView, error)
	RemoveRoomMember(ctx context.Context, roomID, userID ulid.ULID) (*RoomView, error)
}

type gatewayImpl struct {
	repo        repository.Repository
	userGateway usergw.Gateway
}

func New(q sqlc.Querier, userGateway usergw.Gateway) Gateway {
	return &gatewayImpl{
		repo:        repository.New(q),
		userGateway: userGateway,
	}
}
