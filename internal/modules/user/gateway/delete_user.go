package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) DeleteUser(ctx context.Context, id ulid.ULID) error {
	return g.repo.DeleteUser(ctx, id)
}
