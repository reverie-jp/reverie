package gateway

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
)

func (g *gatewayImpl) SearchUsers(ctx context.Context, query string, cursor *time.Time, limit int32) ([]*entity.User, error) {
	return g.repo.SearchUsers(ctx, query, cursor, limit)
}
