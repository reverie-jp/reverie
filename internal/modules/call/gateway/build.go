package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) BuildView(ctx context.Context, requesterID, callID ulid.ULID) (*CallView, error) {
	views, err := g.BuildListViews(ctx, requesterID, []ulid.ULID{callID})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

// BuildListViews materializes calls for the given IDs and their host user
// views in two batched queries, preserving the input order.
func (g *gatewayImpl) BuildListViews(ctx context.Context, requesterID ulid.ULID, callIDs []ulid.ULID) ([]*CallView, error) {
	if len(callIDs) == 0 {
		return []*CallView{}, nil
	}

	calls, err := g.repo.ListCallsByIDs(ctx, callIDs)
	if err != nil {
		return nil, err
	}

	hostIDs := make([]ulid.ULID, 0, len(calls))
	seen := make(map[ulid.ULID]bool, len(calls))
	for _, c := range calls {
		if !seen[c.HostUserID] {
			hostIDs = append(hostIDs, c.HostUserID)
			seen[c.HostUserID] = true
		}
	}
	hostViews, err := g.userGateway.BuildListViews(ctx, requesterID, hostIDs)
	if err != nil {
		return nil, err
	}
	hostByID := make(map[ulid.ULID]*usergw.UserView, len(hostViews))
	for _, v := range hostViews {
		if v != nil && v.User != nil {
			hostByID[v.User.ID] = v
		}
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
