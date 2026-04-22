package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListFollowingIDs(ctx context.Context, followerID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error) {
	return g.repo.ListFollowingIDs(ctx, followerID, cursorID, pageSize)
}
