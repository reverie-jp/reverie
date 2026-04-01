package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
)

type Repository interface {
	GetAuthProviderByProvider(ctx context.Context, provider string, providerUserID string) (*entity.AuthProvider, error)
	CreateAuthProvider(ctx context.Context, params CreateAuthProviderParams) error
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}
