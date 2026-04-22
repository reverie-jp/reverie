package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListNotificationsByRecipient(ctx context.Context, recipientID ulid.ULID, cursorID string, pageSize int32) ([]*entity.Notification, error) {
	rows, err := r.q.ListNotificationsByRecipient(ctx, sqlc.ListNotificationsByRecipientParams{
		RecipientUserID: recipientID,
		CursorID:        cursorID,
		PageSize:        pageSize,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.Notification, 0, len(rows))
	for i := range rows {
		out = append(out, mapper.ToNotification(&rows[i]))
	}
	return out, nil
}
