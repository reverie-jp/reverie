package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) CountUnread(ctx context.Context, recipientID ulid.ULID) (int32, error) {
	return g.repo.CountUnreadNotifications(ctx, recipientID)
}
