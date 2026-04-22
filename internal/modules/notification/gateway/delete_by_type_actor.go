package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) DeleteByTypeActor(ctx context.Context, recipientID ulid.ULID, notifType entity.NotificationType, actorID ulid.ULID) error {
	_, err := g.repo.DeleteNotificationsByTypeActor(ctx, recipientID, notifType, actorID)
	return err
}
