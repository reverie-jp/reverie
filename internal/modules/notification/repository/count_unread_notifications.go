package repository

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) CountUnreadNotifications(ctx context.Context, recipientID ulid.ULID) (int32, error) {
	return r.q.CountUnreadNotifications(ctx, recipientID)
}
