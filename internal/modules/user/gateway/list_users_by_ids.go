package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error) {
	return g.repo.ListUsersByIDs(ctx, ids)
}
