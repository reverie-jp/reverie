package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) MarkNotificationsRead(ctx context.Context, recipientID ulid.ULID, ids []ulid.ULID) (int64, error) {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = id.String()
	}
	return r.q.MarkNotificationsRead(ctx, sqlc.MarkNotificationsReadParams{
		RecipientUserID: recipientID,
		Ids:             strIDs,
	})
}

func (r *RepositoryImpl) MarkAllNotificationsRead(ctx context.Context, recipientID ulid.ULID) (int64, error) {
	return r.q.MarkAllNotificationsRead(ctx, recipientID)
}
