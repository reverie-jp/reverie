package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type Repository interface {
	GetUserByID(ctx context.Context, id ulid.ULID) (*entity.User, error)
	CreateUser(ctx context.Context, params CreateUserParams) error
	DeleteUser(ctx context.Context, id ulid.ULID) error
	GetAuthProviderByProvider(ctx context.Context, provider string, providerUserID string) (*entity.AuthProvider, error)
	CreateAuthProvider(ctx context.Context, params CreateAuthProviderParams) error
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}
