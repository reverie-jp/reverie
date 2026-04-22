package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListFollowerIDs(ctx context.Context, followeeID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error) {
	return g.repo.ListFollowerIDs(ctx, followeeID, cursorID, pageSize)
}
