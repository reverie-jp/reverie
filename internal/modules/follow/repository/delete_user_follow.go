package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) DeleteUserFollow(ctx context.Context, followerID, followeeID ulid.ULID) error {
	return r.q.DeleteUserFollow(ctx, sqlc.DeleteUserFollowParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
}
