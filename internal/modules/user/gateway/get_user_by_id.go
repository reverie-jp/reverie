package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) GetUserByID(ctx context.Context, id ulid.ULID) (*entity.User, error) {
	return g.repo.GetUserByID(ctx, id)
}
