package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type CreateAuthProviderParams struct {
	UserID         ulid.ULID
	Provider       string
	ProviderUserID string
}

func (r *RepositoryImpl) CreateAuthProvider(ctx context.Context, params CreateAuthProviderParams) error {
	return r.q.CreateAuthProvider(ctx, sqlc.CreateAuthProviderParams{
		ID:             ulid.New(),
		UserID:         params.UserID,
		Provider:       sqlc.AuthProvider(params.Provider),
		ProviderUserID: params.ProviderUserID,
	})
}
