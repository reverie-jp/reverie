package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
)

func (g *gatewayImpl) GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error) {
	return g.repo.GetUserByCustomID(ctx, customID)
}
