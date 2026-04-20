package handler

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	chatv1 "reverie.jp/reverie/internal/gen/pb/chat/v1"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/chat/usecase"
)

func toProtoRoom(out *usecase.RoomOutput) *chatv1.Room {
	r := &chatv1.Room{
		Id:                out.ID.String(),
		RoomType:          out.RoomType,
		Name:              out.Name,
		LastMessageText:   out.LastMessageText,
		UnreadCount:       int32(out.UnreadCount),
		IsLastMessageMine: out.IsLastMessageMine,
		IsPinned:          out.IsPinned,
		IsMuted:           out.IsMuted,
	}
	if out.LastMessageAt != nil {
		r.LastMessageAt = timestamppb.New(*out.LastMessageAt)
	}
	if out.OtherUser != nil {
		bio := out.OtherUser.Biography
		r.OtherUser = &userv1.User{
			Id:             out.OtherUser.ID.String(),
			CustomId:       out.OtherUser.CustomID,
			DisplayName:    out.OtherUser.DisplayName,
			Biography:      &bio,
			IsPrivate:      out.OtherUser.IsPrivate,
			FollowerCount:  int32(out.OtherUser.FollowerCount),
			FollowingCount: int32(out.OtherUser.FollowingCount),
			CreateTime:     timestamppb.New(out.OtherUser.CreateTime),
		}
	}
	for _, m := range out.Members {
		bio := m.Biography
		r.Members = append(r.Members, &userv1.User{
			Id:          m.ID.String(),
			CustomId:    m.CustomID,
			DisplayName: m.DisplayName,
			Biography:   &bio,
			IsPrivate:   m.IsPrivate,
			CreateTime:  timestamppb.New(m.CreateTime),
		})
	}
	return r
}

func toProtoMessage(out *usecase.MessageOutput) *chatv1.ChatMessage {
	m := &chatv1.ChatMessage{
		Id:         out.ID.String(),
		Content:    out.Content,
		IsMine:     out.IsMine,
		CreateTime: timestamppb.New(out.CreateTime),
	}
	if out.Sender != nil {
		m.Sender = &userv1.User{
			Id:          out.Sender.ID.String(),
			CustomId:    out.Sender.CustomID,
			DisplayName: out.Sender.DisplayName,
			CreateTime:  timestamppb.New(time.Time{}),
		}
	}
	for _, r := range out.Reactions {
		m.Reactions = append(m.Reactions, &chatv1.MessageReaction{
			Emoji:  r.Emoji,
			Count:  int32(r.Count),
			IsMine: r.IsMine,
		})
	}
	return m
}
