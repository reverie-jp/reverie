package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) FollowUser(ctx context.Context, followerID, followedID ulid.ULID) error {
	return r.q.CreateUserFollow(ctx, sqlc.CreateUserFollowParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
}

func (r *RepositoryImpl) UnfollowUser(ctx context.Context, followerID, followedID ulid.ULID) error {
	return r.q.DeleteUserFollow(ctx, sqlc.DeleteUserFollowParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
}

func (r *RepositoryImpl) IsFollowing(ctx context.Context, followerID, followedID ulid.ULID) (bool, error) {
	_, err := r.q.GetUserFollow(ctx, sqlc.GetUserFollowParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *RepositoryImpl) CountFollowers(ctx context.Context, userID ulid.ULID) (int64, error) {
	return r.q.CountUserFollowers(ctx, userID)
}

func (r *RepositoryImpl) CountFollowing(ctx context.Context, userID ulid.ULID) (int64, error) {
	return r.q.CountUserFollowing(ctx, userID)
}
