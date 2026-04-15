package repository

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type Repository interface {
	GetUserByID(ctx context.Context, id ulid.ULID) (*entity.User, error)
	GetUserByCustomID(ctx context.Context, customID string) (*entity.User, error)
	ListUsersByIDs(ctx context.Context, ids []ulid.ULID) ([]*entity.User, error)
	CreateUser(ctx context.Context, params CreateUserParams) error
	UpdateUser(ctx context.Context, params UpdateUserParams) (*entity.User, error)
	DeleteUser(ctx context.Context, id ulid.ULID) error
}

type UpdateUserParams struct {
	ID          ulid.ULID
	DisplayName string
	Biography   *string
	IsPrivate   bool
	Birthdate   *time.Time
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}
