package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) CreateFollow(ctx context.Context, followerID, followeeID ulid.ULID) error {
	return g.repo.CreateUserFollow(ctx, followerID, followeeID)
}
