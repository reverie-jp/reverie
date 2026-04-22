package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) BuildUserView(ctx context.Context, requesterID, id ulid.ULID) (*UserView, error) {
	views, err := g.BuildListUserViews(ctx, requesterID, []ulid.ULID{id})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

// BuildListUserViews composes UserView from pure user entities and the requester's
// follow relationship flags. Follow counts live on entity.User itself
// (denormalized columns maintained by the user_follows trigger).
func (g *gatewayImpl) BuildListUserViews(ctx context.Context, requesterID ulid.ULID, ids []ulid.ULID) ([]*UserView, error) {
	if len(ids) == 0 {
		return []*UserView{}, nil
	}
	users, err := g.repo.ListUsersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	rels, err := g.followGateway.RelationshipsByUserIDs(ctx, requesterID, ids)
	if err != nil {
		return nil, err
	}
	views := make([]*UserView, len(users))
	for i, u := range users {
		r := rels[u.ID]
		views[i] = &UserView{
			User:         u,
			IsMe:         !requesterID.IsZero() && u.ID == requesterID,
			IsFollowing:  r.IsFollowing,
			IsFollowedBy: r.IsFollowedBy,
		}
	}
	return views, nil
}
