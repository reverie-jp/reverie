package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) BuildNotificationView(ctx context.Context, recipientID ulid.ULID, n *entity.Notification) (*NotificationView, error) {
	views, err := g.BuildListNotificationViews(ctx, recipientID, []*entity.Notification{n})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

// BuildListNotificationViews composes NotificationView from raw entities and
// actor UserViews in one batched call per actor set.
func (g *gatewayImpl) BuildListNotificationViews(ctx context.Context, recipientID ulid.ULID, notifications []*entity.Notification) ([]*NotificationView, error) {
	if len(notifications) == 0 {
		return []*NotificationView{}, nil
	}
	actorIDs := make([]ulid.ULID, 0, len(notifications))
	seen := make(map[ulid.ULID]struct{}, len(notifications))
	for _, n := range notifications {
		if n == nil || n.ActorUserID == nil {
			continue
		}
		if _, ok := seen[*n.ActorUserID]; ok {
			continue
		}
		seen[*n.ActorUserID] = struct{}{}
		actorIDs = append(actorIDs, *n.ActorUserID)
	}

	actorsByID := map[ulid.ULID]*usergw.UserView{}
	if len(actorIDs) > 0 {
		actors, err := g.userGateway.BuildListUserViews(ctx, recipientID, actorIDs)
		if err != nil {
			return nil, err
		}
		for _, v := range actors {
			if v != nil && v.User != nil {
				actorsByID[v.User.ID] = v
			}
		}
	}

	out := make([]*NotificationView, 0, len(notifications))
	for _, n := range notifications {
		if n == nil {
			continue
		}
		var actor *usergw.UserView
		if n.ActorUserID != nil {
			actor = actorsByID[*n.ActorUserID]
		}
		out = append(out, &NotificationView{Notification: n, Actor: actor})
	}
	return out, nil
}
