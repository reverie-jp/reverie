package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func (r *RepositoryImpl) UpdateUser(ctx context.Context, params UpdateUserParams) (*entity.User, error) {
	row, err := r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:          params.ID,
		DisplayName: params.DisplayName,
		Biography:   params.Biography,
		IsPrivate:   params.IsPrivate,
		Birthdate:   params.Birthdate,
	})
	if err != nil {
		return nil, err
	}

	return mapper.ToUser(&row), nil
}
