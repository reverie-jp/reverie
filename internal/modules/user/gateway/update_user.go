package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/modules/user/repository"
)

func (g *gatewayImpl) UpdateUser(ctx context.Context, params UpdateUserParams) (*entity.User, error) {
	return g.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:          params.ID,
		DisplayName: params.DisplayName,
		Biography:   params.Biography,
		IsPrivate:   params.IsPrivate,
		Birthdate:   params.Birthdate,
	})
}
