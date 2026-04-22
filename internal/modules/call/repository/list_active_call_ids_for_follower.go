package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) ListActiveCallIDsForFollower(ctx context.Context, followerID ulid.ULID, staleSeconds int32, cursorID string, pageSize int32) ([]ulid.ULID, error) {
	return r.q.ListActiveCallIDsForFollower(ctx, sqlc.ListActiveCallIDsForFollowerParams{
		FollowerID:   followerID,
		StaleSeconds: staleSeconds,
		CursorID:     cursorID,
		PageSize:     pageSize,
	})
}
