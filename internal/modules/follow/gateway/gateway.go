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
	RelationshipsByUserIDs(ctx context.Context, requesterID ulid.ULID, targetIDs []ulid.ULID) (map[ulid.ULID]Relationship, error)
}

type gatewayImpl struct {
	repo repository.Repository
}

func New(q sqlc.Querier) Gateway {
	return &gatewayImpl{repo: repository.New(q)}
}

func (g *gatewayImpl) CreateFollow(ctx context.Context, followerID, followeeID ulid.ULID) error {
	return g.repo.CreateUserFollow(ctx, followerID, followeeID)
}

func (g *gatewayImpl) DeleteFollow(ctx context.Context, followerID, followeeID ulid.ULID) error {
	return g.repo.DeleteUserFollow(ctx, followerID, followeeID)
}

func (g *gatewayImpl) IsFollowing(ctx context.Context, followerID, followeeID ulid.ULID) (bool, error) {
	return g.repo.IsFollowing(ctx, followerID, followeeID)
}

func (g *gatewayImpl) ListFollowingIDs(ctx context.Context, followerID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error) {
	return g.repo.ListFollowingIDs(ctx, followerID, cursorID, pageSize)
}

func (g *gatewayImpl) ListFollowerIDs(ctx context.Context, followeeID ulid.ULID, cursorID string, pageSize int32) ([]ulid.ULID, error) {
	return g.repo.ListFollowerIDs(ctx, followeeID, cursorID, pageSize)
}

func (g *gatewayImpl) RelationshipsByUserIDs(ctx context.Context, requesterID ulid.ULID, targetIDs []ulid.ULID) (map[ulid.ULID]Relationship, error) {
	rels := make(map[ulid.ULID]Relationship, len(targetIDs))
	for _, id := range targetIDs {
		rels[id] = Relationship{}
	}
	if requesterID.IsZero() || len(targetIDs) == 0 {
		return rels, nil
	}
	following, err := g.repo.ListFollowingEdges(ctx, requesterID, targetIDs)
	if err != nil {
		return nil, err
	}
	for _, id := range following {
		r := rels[id]
		r.IsFollowing = true
		rels[id] = r
	}
	followers, err := g.repo.ListFollowerEdges(ctx, requesterID, targetIDs)
	if err != nil {
		return nil, err
	}
	for _, id := range followers {
		r := rels[id]
		r.IsFollowedBy = true
		rels[id] = r
	}
	return rels, nil
}
