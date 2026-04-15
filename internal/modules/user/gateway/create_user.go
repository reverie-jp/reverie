package gateway

import (
	"context"

	"reverie.jp/reverie/internal/modules/user/repository"
)

func (g *gatewayImpl) CreateUser(ctx context.Context, params CreateUserParams) error {
	return g.repo.CreateUser(ctx, repository.CreateUserParams{
		ID:          params.ID,
		CustomID:    params.CustomID,
		DisplayName: params.DisplayName,
	})
}
