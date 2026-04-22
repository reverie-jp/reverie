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

// BuildListCallViews materializes calls for the given IDs, their host user
// views, and a per-call active-participants list (for avatar stacks on call
// cards). Preserves input order.
//
// Query budget: 3 batched queries total — calls, active participants per
// call, and users for (hosts ∪ authenticated participants). Guests are
// included in the participant list with a nil User. Frontend derives counts
// from the returned slice length.
func (g *gatewayImpl) BuildListCallViews(ctx context.Context, requesterID ulid.ULID, callIDs []ulid.ULID) ([]*CallView, error) {
	if len(callIDs) == 0 {
		return []*CallView{}, nil
	}

	calls, err := g.repo.ListCallsByIDs(ctx, callIDs)
	if err != nil {
		return nil, err
	}

	participantsByCall, err := g.repo.ListActiveParticipantsByCallIDs(ctx, callIDs, entity.ParticipantStaleSeconds)
	if err != nil {
		return nil, err
	}

	// Batch-fetch every user we'll need (hosts + authenticated participants)
	// in one call to the user gateway, then slice per-call.
	userIDSet := make(map[ulid.ULID]struct{}, len(calls)*2)
	for _, c := range calls {
		userIDSet[c.HostUserID] = struct{}{}
	}
	for _, ps := range participantsByCall {
		for _, p := range ps {
			if p.UserID != nil {
				userIDSet[*p.UserID] = struct{}{}
			}
		}
	}
	userIDs := make([]ulid.ULID, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}
	userByID, err := g.buildUserViewMap(ctx, requesterID, userIDs)
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
		ps := participantsByCall[id]
		pvs := make([]*CallParticipantView, 0, len(ps))
		for _, p := range ps {
			var u *usergw.UserView
			if p.UserID != nil {
				u = userByID[*p.UserID]
			}
			pvs = append(pvs, &CallParticipantView{
				Participant:          p,
				User:                 u,
				IsCurrentlyConnected: true,
			})
		}
		out = append(out, &CallView{
			Call:               c,
			Host:               userByID[c.HostUserID],
			ActiveParticipants: pvs,
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
