package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListFollowerIDs(ctx context.Context, followeeID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error) {
	return r.q.ListFollowerIDs(ctx, sqlc.ListFollowerIDsParams{
		FolloweeID: followeeID,
		CursorID:   cursorID,
		PageSize:   pageSize,
	})
}
