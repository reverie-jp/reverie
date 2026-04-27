package usecase

import (
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/modules/chat/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type ChatUserOutput struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
	Biography   string
	IsPrivate   bool
	CreateTime  time.Time
}

type RoomOutput struct {
	ID                ulid.ULID
	RoomType          string
	Name              string
	OtherUser         *ChatUserOutput
	Members           []*ChatUserOutput
	LastMessageText   string
	LastMessageAt     *time.Time
	UnreadCount       int64
	IsLastMessageMine bool
	IsPinned          bool
	IsMuted           bool
}

type ReactionOutput struct {
	Emoji  string
	Count  int64
	IsMine bool
}

type MessageOutput struct {
	ID         ulid.ULID
	Content    string
	Sender     *entity.User
	IsMine     bool
	CreateTime time.Time
	Reactions  []*ReactionOutput
}

func toRoomOutput(view *gateway.RoomView, requestorID ulid.ULID) *RoomOutput {
	out := &RoomOutput{
		ID:            view.Room.ID,
		RoomType:      view.Room.RoomType,
		LastMessageAt: view.LastMessageAt,
		UnreadCount:   view.UnreadCount,
		IsPinned:      view.IsPinned,
		IsMuted:       view.IsMuted,
	}
	if view.Room.Name != nil {
		out.Name = *view.Room.Name
	}
	if view.LastMessageText != nil {
		out.LastMessageText = *view.LastMessageText
	}
	if view.LastMessageSendID != nil {
		out.IsLastMessageMine = *view.LastMessageSendID == requestorID
	}
	if view.OtherUser != nil {
		out.OtherUser = entityUserToOutput(view.OtherUser)
	}
	for _, m := range view.Members {
		out.Members = append(out.Members, entityUserToOutput(m))
	}
	return out
}

func entityUserToOutput(u *entity.User) *ChatUserOutput {
	return &ChatUserOutput{
		ID:          u.ID,
		CustomID:    u.CustomID,
		DisplayName: u.DisplayName,
		Biography:   u.Biography,
		IsPrivate:   u.IsPrivate,
		CreateTime:  u.CreateTime,
	}
}

func toMessageOutput(view *gateway.MessageView, requestorID ulid.ULID) *MessageOutput {
	content := ""
	if view.Message.Content != nil {
		content = *view.Message.Content
	}
	reactions := make([]*ReactionOutput, len(view.Reactions))
	for i, r := range view.Reactions {
		reactions[i] = &ReactionOutput{Emoji: r.Emoji, Count: r.Count, IsMine: r.IsMine}
	}
	return &MessageOutput{
		ID:         view.Message.ID,
		Content:    content,
		Sender:     view.Sender,
		IsMine:     view.Message.SenderID == requestorID,
		CreateTime: view.Message.CreateTime,
		Reactions:  reactions,
	}
}
