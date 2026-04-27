package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) BuildUserView(ctx context.Context, requesterID ulid.ULID, id ulid.ULID) (*UserView, error) {
	views, err := g.BuildListUserViews(ctx, requesterID, []ulid.ULID{id})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	view := views[0]

	followerCount, err := g.repo.CountFollowers(ctx, id)
	if err != nil {
		return nil, err
	}
	followingCount, err := g.repo.CountFollowing(ctx, id)
	if err != nil {
		return nil, err
	}
	view.FollowerCount = followerCount
	view.FollowingCount = followingCount

	return view, nil
}

func (g *gatewayImpl) BuildListUserViews(ctx context.Context, requesterID ulid.ULID, ids []ulid.ULID) ([]*UserView, error) {
	if len(ids) == 0 {
		return []*UserView{}, nil
	}

	users, err := g.repo.ListUsersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = id.String()
	}

	followingSet := map[ulid.ULID]bool{}
	followerSet := map[ulid.ULID]bool{}

	if !requesterID.IsZero() {
		followingIDs, err := g.repo.ListFollowingEdges(ctx, requesterID, idStrs)
		if err != nil {
			return nil, err
		}
		for _, id := range followingIDs {
			followingSet[id] = true
		}

		followerIDs, err := g.repo.ListFollowerEdges(ctx, requesterID, idStrs)
		if err != nil {
			return nil, err
		}
		for _, id := range followerIDs {
			followerSet[id] = true
		}
	}

	views := make([]*UserView, len(users))
	for i, u := range users {
		views[i] = &UserView{
			User:         u,
			IsMe:         u.ID == requesterID,
			IsFollowing:  followingSet[u.ID],
			IsFollowedBy: followerSet[u.ID],
		}
	}
	return views, nil
}
