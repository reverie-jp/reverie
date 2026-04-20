package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/modules/chat/repository"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (g *gatewayImpl) ListRooms(ctx context.Context, userID ulid.ULID) ([]*RoomView, error) {
	rows, err := g.repo.ListRoomsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	views := make([]*RoomView, 0, len(rows))
	for _, row := range rows {
		view, err := g.buildRoomView(ctx, row, userID)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

func (g *gatewayImpl) CreateGroupRoom(ctx context.Context, creatorID ulid.ULID, name string, memberIDs []ulid.ULID) (*RoomView, error) {
	roomID := ulid.New()
	_, err := g.repo.CreateRoom(ctx, roomID, "group", &name)
	if err != nil {
		return nil, err
	}
	allMembers := append([]ulid.ULID{creatorID}, memberIDs...)
	seen := map[ulid.ULID]bool{}
	for _, mid := range allMembers {
		if seen[mid] {
			continue
		}
		seen[mid] = true
		if err := g.repo.CreateRoomMember(ctx, ulid.New(), roomID, mid); err != nil {
			return nil, err
		}
	}
	return g.GetRoomView(ctx, roomID, creatorID)
}

func (g *gatewayImpl) UpdateRoom(ctx context.Context, roomID ulid.ULID, name string) (*RoomView, error) {
	if _, err := g.repo.UpdateRoomName(ctx, roomID, &name); err != nil {
		return nil, err
	}
	return g.getGroupRoomView(ctx, roomID)
}

func (g *gatewayImpl) ListRoomMembers(ctx context.Context, roomID ulid.ULID) ([]*entity.User, error) {
	memberIDs, err := g.repo.GetRoomMembers(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if len(memberIDs) == 0 {
		return nil, nil
	}
	return g.userGateway.ListUsersByIDs(ctx, memberIDs)
}

func (g *gatewayImpl) AddRoomMember(ctx context.Context, roomID, userID ulid.ULID) (*RoomView, error) {
	if err := g.repo.CreateRoomMember(ctx, ulid.New(), roomID, userID); err != nil {
		return nil, err
	}
	return g.getGroupRoomView(ctx, roomID)
}

func (g *gatewayImpl) RemoveRoomMember(ctx context.Context, roomID, userID ulid.ULID) (*RoomView, error) {
	if err := g.repo.RemoveRoomMember(ctx, roomID, userID); err != nil {
		return nil, err
	}
	return g.getGroupRoomView(ctx, roomID)
}

func (g *gatewayImpl) getGroupRoomView(ctx context.Context, roomID ulid.ULID) (*RoomView, error) {
	room, err := g.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, xerrors.ErrNotFound
		}
		return nil, err
	}
	members, err := g.ListRoomMembers(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return &RoomView{Room: room, Members: members}, nil
}

func (g *gatewayImpl) buildRoomView(ctx context.Context, row *repository.RoomRow, userID ulid.ULID) (*RoomView, error) {
	room := &entity.Room{
		ID:       row.ID,
		RoomType: row.RoomType,
		Name:     row.Name,
	}

	view := &RoomView{
		Room:            room,
		LastMessageText: row.LastMessageText,
		LastMessageAt:   row.LastMessageAt,
		UnreadCount:     row.UnreadCount,
		IsPinned:        row.IsPinned,
		IsMuted:         row.IsMuted,
	}

	if row.LastMessageSenderID != nil {
		id, err := ulid.Parse(*row.LastMessageSenderID)
		if err == nil {
			view.LastMessageSendID = &id
		}
	}

	if row.RoomType == "direct" {
		otherMemberID, err := g.repo.GetRoomOtherMember(ctx, row.ID, userID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		if err == nil {
			otherUser, err := g.userGateway.GetUserByID(ctx, otherMemberID)
			if err != nil {
				return nil, err
			}
			view.OtherUser = otherUser
		}
	} else if row.RoomType == "group" {
		memberIDs, err := g.repo.GetRoomMembers(ctx, row.ID)
		if err == nil && len(memberIDs) > 0 {
			members, err := g.userGateway.ListUsersByIDs(ctx, memberIDs)
			if err == nil {
				view.Members = members
			}
		}
	}

	return view, nil
}

func (g *gatewayImpl) GetOrCreateDirectRoom(ctx context.Context, myID, otherID ulid.ULID) (*RoomView, error) {
	existingID, err := g.repo.GetDirectRoom(ctx, myID, otherID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	var roomID ulid.ULID
	if errors.Is(err, pgx.ErrNoRows) {
		newRoomID := ulid.New()
		_, err = g.repo.CreateRoom(ctx, newRoomID, "direct", nil)
		if err != nil {
			return nil, err
		}
		if err := g.repo.CreateRoomMember(ctx, ulid.New(), newRoomID, myID); err != nil {
			return nil, err
		}
		if err := g.repo.CreateRoomMember(ctx, ulid.New(), newRoomID, otherID); err != nil {
			return nil, err
		}
		roomID = newRoomID
	} else {
		roomID = existingID
	}

	return g.GetRoomView(ctx, roomID, myID)
}

func (g *gatewayImpl) GetRoomView(ctx context.Context, roomID, userID ulid.ULID) (*RoomView, error) {
	rooms, err := g.ListRooms(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, r := range rooms {
		if r.Room.ID == roomID {
			return r, nil
		}
	}

	room, err := g.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, xerrors.ErrNotFound
		}
		return nil, err
	}

	view := &RoomView{Room: room}

	if room.RoomType == "direct" {
		otherMemberID, err := g.repo.GetRoomOtherMember(ctx, roomID, userID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		if err == nil {
			otherUser, err := g.userGateway.GetUserByID(ctx, otherMemberID)
			if err != nil {
				return nil, err
			}
			view.OtherUser = otherUser
		}
	}

	return view, nil
}

func (g *gatewayImpl) ListMessages(ctx context.Context, roomID ulid.ULID, cursor *time.Time, limit int32, requestorID ulid.ULID) ([]*MessageView, error) {
	return g.listMessagesWithReactions(ctx, roomID, cursor, limit, requestorID)
}

func (g *gatewayImpl) listMessagesWithReactions(ctx context.Context, roomID ulid.ULID, cursor *time.Time, limit int32, requestorID ulid.ULID) ([]*MessageView, error) {
	msgs, err := g.repo.ListMessages(ctx, roomID, cursor, limit)
	if err != nil {
		return nil, err
	}

	senderIDs := make([]ulid.ULID, 0, len(msgs))
	seen := map[ulid.ULID]bool{}
	for _, m := range msgs {
		if !seen[m.SenderID] {
			senderIDs = append(senderIDs, m.SenderID)
			seen[m.SenderID] = true
		}
	}

	senderMap := map[ulid.ULID]*entity.User{}
	if len(senderIDs) > 0 {
		users, err := g.userGateway.ListUsersByIDs(ctx, senderIDs)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			senderMap[u.ID] = u
		}
	}

	reactionMap := map[ulid.ULID][]*ReactionView{}
	if len(msgs) > 0 {
		msgIDs := make([]ulid.ULID, len(msgs))
		for i, m := range msgs {
			msgIDs[i] = m.ID
		}
		reactionRows, err := g.repo.ListReactionsByMessages(ctx, msgIDs, requestorID)
		if err != nil {
			return nil, err
		}
		for _, r := range reactionRows {
			reactionMap[r.MessageID] = append(reactionMap[r.MessageID], &ReactionView{
				Emoji:  r.Emoji,
				Count:  r.Count,
				IsMine: r.IsMine,
			})
		}
	}

	views := make([]*MessageView, len(msgs))
	for i, m := range msgs {
		views[i] = &MessageView{
			Message:   m,
			Sender:    senderMap[m.SenderID],
			Reactions: reactionMap[m.ID],
		}
	}
	return views, nil
}

func (g *gatewayImpl) SendMessage(ctx context.Context, roomID, senderID ulid.ULID, content string) (*MessageView, error) {
	msg, err := g.repo.CreateMessage(ctx, ulid.New(), roomID, senderID, content)
	if err != nil {
		return nil, err
	}
	sender, err := g.userGateway.GetUserByID(ctx, senderID)
	if err != nil {
		return nil, err
	}
	return &MessageView{Message: msg, Sender: sender}, nil
}

func (g *gatewayImpl) MarkRoomAsRead(ctx context.Context, roomID, userID ulid.ULID) error {
	return g.repo.MarkRoomAsRead(ctx, roomID, userID)
}

func (g *gatewayImpl) PinRoom(ctx context.Context, userID, roomID ulid.ULID) (*RoomView, error) {
	if err := g.repo.PinRoom(ctx, ulid.New(), userID, roomID); err != nil {
		return nil, err
	}
	return g.GetRoomView(ctx, roomID, userID)
}

func (g *gatewayImpl) UnpinRoom(ctx context.Context, userID, roomID ulid.ULID) (*RoomView, error) {
	if err := g.repo.UnpinRoom(ctx, userID, roomID); err != nil {
		return nil, err
	}
	return g.GetRoomView(ctx, roomID, userID)
}

func (g *gatewayImpl) MuteRoom(ctx context.Context, userID, roomID ulid.ULID) (*RoomView, error) {
	if err := g.repo.MuteRoom(ctx, roomID, userID, true); err != nil {
		return nil, err
	}
	return g.GetRoomView(ctx, roomID, userID)
}

func (g *gatewayImpl) UnmuteRoom(ctx context.Context, userID, roomID ulid.ULID) (*RoomView, error) {
	if err := g.repo.MuteRoom(ctx, roomID, userID, false); err != nil {
		return nil, err
	}
	return g.GetRoomView(ctx, roomID, userID)
}

func (g *gatewayImpl) LeaveRoom(ctx context.Context, roomID, userID ulid.ULID) error {
	return g.repo.LeaveRoom(ctx, roomID, userID)
}

func (g *gatewayImpl) AddReaction(ctx context.Context, messageID, userID ulid.ULID, emoji string) (*MessageView, error) {
	if err := g.repo.AddReaction(ctx, messageID, userID, emoji); err != nil {
		return nil, err
	}
	return g.buildMessageView(ctx, messageID, userID)
}

func (g *gatewayImpl) RemoveReaction(ctx context.Context, messageID, userID ulid.ULID, emoji string) (*MessageView, error) {
	if err := g.repo.RemoveReaction(ctx, messageID, userID, emoji); err != nil {
		return nil, err
	}
	return g.buildMessageView(ctx, messageID, userID)
}

func (g *gatewayImpl) buildMessageView(ctx context.Context, messageID, requestorID ulid.ULID) (*MessageView, error) {
	msg, err := g.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	sender, err := g.userGateway.GetUserByID(ctx, msg.SenderID)
	if err != nil {
		return nil, err
	}
	reactionRows, err := g.repo.ListReactionsByMessage(ctx, messageID, requestorID)
	if err != nil {
		return nil, err
	}
	reactions := make([]*ReactionView, len(reactionRows))
	for i, r := range reactionRows {
		reactions[i] = &ReactionView{Emoji: r.Emoji, Count: r.Count, IsMine: r.IsMine}
	}
	return &MessageView{
		Message:   msg,
		Sender:    sender,
		Reactions: reactions,
	}, nil
}
