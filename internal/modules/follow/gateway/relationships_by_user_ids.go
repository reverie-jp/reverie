package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

// RelationshipsByUserIDs returns a Relationship for every target ID regardless
// of whether an edge exists, so callers can do a single map lookup per user.
// Guest requesters (zero ULID) short-circuit to all-false.
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
