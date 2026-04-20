package gateway

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListFollowing(ctx context.Context, followerID ulid.ULID, cursor *time.Time, limit int32) ([]*entity.User, error) {
	return g.repo.ListFollowing(ctx, followerID, cursor, limit)
}
