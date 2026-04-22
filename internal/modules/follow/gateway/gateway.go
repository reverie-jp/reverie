package gateway

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/follow/repository"
	"reverie.jp/reverie/internal/platform/ulid"
)

type Relationship struct {
	IsFollowing  bool
	IsFollowedBy bool
}

type Gateway interface {
	CreateFollow(ctx context.Context, followerID, followeeID ulid.ULID) error
	DeleteFollow(ctx context.Context, followerID, followeeID ulid.ULID) error
	IsFollowing(ctx context.Context, followerID, followeeID ulid.ULID) (bool, error)
	ListFollowingIDs(ctx context.Context, followerID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error)
	ListFollowerIDs(ctx context.Context, followeeID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error)
	// ListAllFollowerIDs returns every follower of followeeID in a single
	// query. Used by fan-out writers (call creation) that touch all followers.
	ListAllFollowerIDs(ctx context.Context, followeeID ulid.ULID) ([]ulid.ULID, error)
	RelationshipsByUserIDs(ctx context.Context, requesterID ulid.ULID, targetIDs []ulid.ULID) (map[ulid.ULID]Relationship, error)
}

type gatewayImpl struct {
	repo repository.Repository
}

func New(q sqlc.Querier) Gateway {
	return &gatewayImpl{repo: repository.New(q)}
}
