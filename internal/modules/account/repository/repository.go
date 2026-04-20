package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type Repository interface {
	GetAuthProviderByProvider(ctx context.Context, provider string, providerUserID string) (*entity.AuthProvider, error)
	CreateAuthProvider(ctx context.Context, params CreateAuthProviderParams) error
	CreateRefreshToken(ctx context.Context, params CreateRefreshTokenParams) error
	GetRefreshTokenByRaw(ctx context.Context, raw string) (*entity.RefreshToken, error)
	DeleteRefreshTokenByRaw(ctx context.Context, raw string, userID ulid.ULID) error
	DeleteExpiredRefreshTokensByUserID(ctx context.Context, userID ulid.ULID) error
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}
