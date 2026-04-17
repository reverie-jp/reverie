package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) FollowUser(ctx context.Context, followerID, followedID ulid.ULID) error {
	return g.repo.FollowUser(ctx, followerID, followedID)
}

func (g *gatewayImpl) UnfollowUser(ctx context.Context, followerID, followedID ulid.ULID) error {
	return g.repo.UnfollowUser(ctx, followerID, followedID)
}

func (g *gatewayImpl) IsFollowing(ctx context.Context, followerID, followedID ulid.ULID) (bool, error) {
	return g.repo.IsFollowing(ctx, followerID, followedID)
}

func (g *gatewayImpl) CountFollowers(ctx context.Context, userID ulid.ULID) (int64, error) {
	return g.repo.CountFollowers(ctx, userID)
}

func (g *gatewayImpl) CountFollowing(ctx context.Context, userID ulid.ULID) (int64, error) {
	return g.repo.CountFollowing(ctx, userID)
}
