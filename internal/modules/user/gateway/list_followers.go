package gateway

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListFollowers(ctx context.Context, followedID ulid.ULID, cursor *time.Time, limit int32) ([]*entity.User, error) {
	return g.repo.ListFollowers(ctx, followedID, cursor, limit)
}
