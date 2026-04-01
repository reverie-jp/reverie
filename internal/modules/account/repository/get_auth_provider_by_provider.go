package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func (r *RepositoryImpl) GetAuthProviderByProvider(ctx context.Context, provider string, providerUserID string) (*entity.AuthProvider, error) {
	row, err := r.q.GetAuthProviderByProvider(ctx, sqlc.GetAuthProviderByProviderParams{
		Provider:       sqlc.AuthProvider(provider),
		ProviderUserID: providerUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.AuthProvider{
		UserID: row.UserID,
	}, nil
}
