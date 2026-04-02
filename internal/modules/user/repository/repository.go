package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type Repository interface {
	ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error)
	GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error)
	CreateUser(ctx context.Context, params CreateUserParams) error
	DeleteUser(ctx context.Context, id ulid.ULID) error
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}
