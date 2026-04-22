package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListFollowingIDs(ctx context.Context, followerID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error) {
	return r.q.ListFollowingIDs(ctx, sqlc.ListFollowingIDsParams{
		FollowerID: followerID,
		CursorID:   cursorID,
		PageSize:   pageSize,
	})
}
