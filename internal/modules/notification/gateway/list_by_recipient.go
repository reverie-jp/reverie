package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) ListByRecipient(ctx context.Context, recipientID ulid.ULID, cursorID string, pageSize int32) ([]*NotificationView, error) {
	notifs, err := g.repo.ListNotificationsByRecipient(ctx, recipientID, cursorID, pageSize)
	if err != nil {
		return nil, err
	}
	return g.BuildListNotificationViews(ctx, recipientID, notifs)
}
