package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type Repository interface {
	CreateUserFollow(ctx context.Context, followerID, followeeID ulid.ULID) error
	DeleteUserFollow(ctx context.Context, followerID, followeeID ulid.ULID) error
	IsFollowing(ctx context.Context, followerID, followeeID ulid.ULID) (bool, error)
	ListFollowingIDs(ctx context.Context, followerID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error)
	ListFollowerIDs(ctx context.Context, followeeID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error)
	ListFollowingEdges(ctx context.Context, followerID ulid.ULID, targetIDs []ulid.ULID) ([]ulid.ULID, error)
	ListFollowerEdges(ctx context.Context, followeeID ulid.ULID, targetIDs []ulid.ULID) ([]ulid.ULID, error)
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}
