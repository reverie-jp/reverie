package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) MarkRead(ctx context.Context, recipientID ulid.ULID, ids []ulid.ULID) (int32, error) {
	n, err := g.repo.MarkNotificationsRead(ctx, recipientID, ids)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}

func (g *gatewayImpl) MarkAllRead(ctx context.Context, recipientID ulid.ULID) (int32, error) {
	n, err := g.repo.MarkAllNotificationsRead(ctx, recipientID)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}
