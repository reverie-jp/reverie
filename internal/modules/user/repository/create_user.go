package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type CreateUserParams struct {
	ID          ulid.ULID
	CustomID    string
	DisplayName string
}

func (r *RepositoryImpl) CreateUser(ctx context.Context, params CreateUserParams) error {
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:          params.ID,
		CustomID:    params.CustomID,
		DisplayName: params.DisplayName,
	})
}
