package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) CreateUserFollow(ctx context.Context, followerID, followeeID ulid.ULID) error {
	return r.q.CreateUserFollow(ctx, sqlc.CreateUserFollowParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
}
