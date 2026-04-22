package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) IsFollowing(ctx context.Context, followerID, followeeID ulid.ULID) (bool, error) {
	return r.q.IsFollowing(ctx, sqlc.IsFollowingParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
}
