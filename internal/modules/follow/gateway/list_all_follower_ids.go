package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListAllFollowerIDs(ctx context.Context, followeeID ulid.ULID) ([]ulid.ULID, error) {
	return g.repo.ListAllFollowerIDs(ctx, followeeID)
}
