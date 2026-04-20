package repository

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *repositoryImpl) ListRoomsByUser(ctx context.Context, userID ulid.ULID) ([]*RoomRow, error) {
	rows, err := r.q.ListRoomsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*RoomRow, len(rows))
	for i, row := range rows {
		var senderID *string
		if row.LastMessageSenderID != nil {
			s := *row.LastMessageSenderID
			senderID = &s
		}
		result[i] = &RoomRow{
			ID:                  row.ID,
			RoomType:            string(row.RoomType),
			Name:                row.Name,
			IsMuted:             row.IsMuted,
			LastReadAt:          row.LastReadAt,
			IsPinned:            row.IsPinned,
			LastMessageText:     row.LastMessageText,
			LastMessageSenderID: senderID,
			LastMessageAt:       row.LastMessageAt,
			UnreadCount:         row.UnreadCount,
		}
	}
	return result, nil
}

func (r *repositoryImpl) GetRoomByID(ctx context.Context, id ulid.ULID) (*entity.Room, error) {
	row, err := r.q.GetRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Room{
		ID:         row.ID,
		RoomType:   string(row.RoomType),
		Name:       row.Name,
		CreateTime: row.CreateTime,
		UpdateTime: row.UpdateTime,
	}, nil
}

func (r *repositoryImpl) GetDirectRoom(ctx context.Context, userID1, userID2 ulid.ULID) (ulid.ULID, error) {
	return r.q.GetDirectRoom(ctx, userID1, userID2)
}

func (r *repositoryImpl) CreateRoom(ctx context.Context, id ulid.ULID, roomType string, name *string) (*entity.Room, error) {
	row, err := r.q.CreateRoom(ctx, sqlc.CreateRoomParams{
		ID:       id,
		RoomType: roomType,
		Name:     name,
	})
	if err != nil {
		return nil, err
	}
	return &entity.Room{
		ID:         row.ID,
		RoomType:   string(row.RoomType),
		Name:       row.Name,
		CreateTime: row.CreateTime,
		UpdateTime: row.UpdateTime,
	}, nil
}

func (r *repositoryImpl) CreateRoomMember(ctx context.Context, id, roomID, userID ulid.ULID) error {
	return r.q.CreateRoomMember(ctx, sqlc.CreateRoomMemberParams{
		ID:     id,
		RoomID: roomID,
		UserID: userID,
	})
}

func (r *repositoryImpl) GetRoomMember(ctx context.Context, roomID, userID ulid.ULID) (*sqlc.RoomMember, error) {
	m, err := r.q.GetRoomMember(ctx, sqlc.GetRoomMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repositoryImpl) GetRoomOtherMember(ctx context.Context, roomID, userID ulid.ULID) (ulid.ULID, error) {
	return r.q.GetRoomOtherMember(ctx, sqlc.GetRoomOtherMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
}

func (r *repositoryImpl) GetMessageByID(ctx context.Context, id ulid.ULID) (*entity.Message, error) {
	row, err := r.q.GetMessageByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Message{
		ID:         row.ID,
		RoomID:     row.RoomID,
		SenderID:   row.SenderID,
		Content:    row.Content,
		IsDeleted:  row.IsDeleted,
		IsEdited:   row.IsEdited,
		CreateTime: row.CreateTime,
	}, nil
}

func (r *repositoryImpl) ListMessages(ctx context.Context, roomID ulid.ULID, cursor *time.Time, limit int32) ([]*entity.Message, error) {
	rows, err := r.q.ListMessages(ctx, sqlc.ListMessagesParams{
		RoomID: roomID,
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*entity.Message, len(rows))
	for i, row := range rows {
		result[i] = &entity.Message{
			ID:         row.ID,
			RoomID:     row.RoomID,
			SenderID:   row.SenderID,
			Content:    row.Content,
			IsDeleted:  row.IsDeleted,
			IsEdited:   row.IsEdited,
			CreateTime: row.CreateTime,
		}
	}
	return result, nil
}

func (r *repositoryImpl) CreateMessage(ctx context.Context, id, roomID, senderID ulid.ULID, content string) (*entity.Message, error) {
	row, err := r.q.CreateMessage(ctx, sqlc.CreateMessageParams{
		ID:       id,
		RoomID:   roomID,
		SenderID: senderID,
		Content:  &content,
	})
	if err != nil {
		return nil, err
	}
	return &entity.Message{
		ID:         row.ID,
		RoomID:     row.RoomID,
		SenderID:   row.SenderID,
		Content:    row.Content,
		IsDeleted:  row.IsDeleted,
		IsEdited:   row.IsEdited,
		CreateTime: row.CreateTime,
	}, nil
}

func (r *repositoryImpl) MarkRoomAsRead(ctx context.Context, roomID, userID ulid.ULID) error {
	return r.q.MarkRoomAsRead(ctx, sqlc.MarkRoomAsReadParams{
		RoomID: roomID,
		UserID: userID,
	})
}

func (r *repositoryImpl) PinRoom(ctx context.Context, id, userID, roomID ulid.ULID) error {
	return r.q.PinRoom(ctx, sqlc.PinRoomParams{
		ID:     id,
		UserID: userID,
		RoomID: roomID,
	})
}

func (r *repositoryImpl) UnpinRoom(ctx context.Context, userID, roomID ulid.ULID) error {
	return r.q.UnpinRoom(ctx, sqlc.UnpinRoomParams{
		UserID: userID,
		RoomID: roomID,
	})
}

func (r *repositoryImpl) MuteRoom(ctx context.Context, roomID, userID ulid.ULID, muted bool) error {
	return r.q.MuteRoom(ctx, sqlc.MuteRoomParams{
		RoomID:  roomID,
		UserID:  userID,
		IsMuted: muted,
	})
}

func (r *repositoryImpl) LeaveRoom(ctx context.Context, roomID, userID ulid.ULID) error {
	return r.q.LeaveRoom(ctx, sqlc.LeaveRoomParams{
		RoomID: roomID,
		UserID: userID,
	})
}

func (r *repositoryImpl) GetRoomMembers(ctx context.Context, roomID ulid.ULID) ([]ulid.ULID, error) {
	return r.q.GetRoomMembers(ctx, roomID)
}

func (r *repositoryImpl) RemoveRoomMember(ctx context.Context, roomID, userID ulid.ULID) error {
	return r.q.RemoveRoomMember(ctx, sqlc.RemoveRoomMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
}

func (r *repositoryImpl) UpdateRoomName(ctx context.Context, roomID ulid.ULID, name *string) (*entity.Room, error) {
	row, err := r.q.UpdateRoomName(ctx, sqlc.UpdateRoomNameParams{
		ID:   roomID,
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &entity.Room{
		ID:         row.ID,
		RoomType:   string(row.RoomType),
		Name:       row.Name,
		CreateTime: row.CreateTime,
		UpdateTime: row.UpdateTime,
	}, nil
}
