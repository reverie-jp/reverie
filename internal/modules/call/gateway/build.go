package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) BuildCallView(ctx context.Context, requesterID, callID ulid.ULID) (*CallView, error) {
	views, err := g.BuildListCallViews(ctx, requesterID, []ulid.ULID{callID})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

// BuildListCallViews materializes calls for the given IDs and their host user
// views in two batched queries, preserving the input order.
func (g *gatewayImpl) BuildListCallViews(ctx context.Context, requesterID ulid.ULID, callIDs []ulid.ULID) ([]*CallView, error) {
	if len(callIDs) == 0 {
		return []*CallView{}, nil
	}

	calls, err := g.repo.ListCallsByIDs(ctx, callIDs)
	if err != nil {
		return nil, err
	}

	hostIDs := uniqueIDs(calls, func(c *entity.Call) ulid.ULID { return c.HostUserID })
	hostByID, err := g.buildUserViewMap(ctx, requesterID, hostIDs)
	if err != nil {
		return nil, err
	}

	callByID := make(map[ulid.ULID]*entity.Call, len(calls))
	for _, c := range calls {
		callByID[c.ID] = c
	}

	out := make([]*CallView, 0, len(callIDs))
	for _, id := range callIDs {
		c, ok := callByID[id]
		if !ok {
			continue
		}
		out = append(out, &CallView{
			Call: c,
			Host: hostByID[c.HostUserID],
		})
	}
	return out, nil
}

// BuildListParticipantViews wraps participants with their user view (if
// authenticated) and derives IsCurrentlyConnected from heartbeat state.
func (g *gatewayImpl) BuildListParticipantViews(ctx context.Context, requesterID ulid.ULID, participants []*entity.CallParticipant) ([]*CallParticipantView, error) {
	if len(participants) == 0 {
		return []*CallParticipantView{}, nil
	}

	userIDs := make([]ulid.ULID, 0, len(participants))
	seen := make(map[ulid.ULID]bool, len(participants))
	for _, p := range participants {
		if p.UserID != nil && !seen[*p.UserID] {
			userIDs = append(userIDs, *p.UserID)
			seen[*p.UserID] = true
		}
	}
	userByID, err := g.buildUserViewMap(ctx, requesterID, userIDs)
	if err != nil {
		return nil, err
	}

	out := make([]*CallParticipantView, len(participants))
	for i, p := range participants {
		var user *usergw.UserView
		if p.UserID != nil {
			user = userByID[*p.UserID]
		}
		out[i] = &CallParticipantView{
			Participant:          p,
			User:                 user,
			IsCurrentlyConnected: p.IsCurrentlyConnected(),
		}
	}
	return out, nil
}

func (g *gatewayImpl) BuildListCallBanViews(ctx context.Context, requesterID ulid.ULID, bans []*entity.CallBan) ([]*CallBanView, error) {
	if len(bans) == 0 {
		return []*CallBanView{}, nil
	}
	userIDs := uniqueIDs(bans, func(b *entity.CallBan) ulid.ULID { return b.UserID })
	userByID, err := g.buildUserViewMap(ctx, requesterID, userIDs)
	if err != nil {
		return nil, err
	}
	out := make([]*CallBanView, len(bans))
	for i, b := range bans {
		out[i] = &CallBanView{Ban: b, User: userByID[b.UserID]}
	}
	return out, nil
}

func (g *gatewayImpl) buildUserViewMap(ctx context.Context, requesterID ulid.ULID, userIDs []ulid.ULID) (map[ulid.ULID]*usergw.UserView, error) {
	views, err := g.userGateway.BuildListUserViews(ctx, requesterID, userIDs)
	if err != nil {
		return nil, err
	}
	byID := make(map[ulid.ULID]*usergw.UserView, len(views))
	for _, v := range views {
		if v != nil && v.User != nil {
			byID[v.User.ID] = v
		}
	}
	return byID, nil
}

func uniqueIDs[T any](items []T, pick func(T) ulid.ULID) []ulid.ULID {
	out := make([]ulid.ULID, 0, len(items))
	seen := make(map[ulid.ULID]bool, len(items))
	for _, it := range items {
		id := pick(it)
		if !seen[id] {
			out = append(out, id)
			seen[id] = true
		}
	}
	return out
}
