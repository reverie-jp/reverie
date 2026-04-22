package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) IsFollowing(ctx context.Context, followerID, followeeID ulid.ULID) (bool, error) {
	return g.repo.IsFollowing(ctx, followerID, followeeID)
}
