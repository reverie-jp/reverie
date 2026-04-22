package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) UpdateLastSeen(ctx context.Context, id ulid.ULID) error {
	return g.repo.UpdateUserLastSeen(ctx, id)
}
