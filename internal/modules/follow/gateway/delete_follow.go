package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) DeleteFollow(ctx context.Context, followerID, followeeID ulid.ULID) error {
	return g.repo.DeleteUserFollow(ctx, followerID, followeeID)
}
