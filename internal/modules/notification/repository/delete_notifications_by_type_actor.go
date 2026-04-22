package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (r *RepositoryImpl) DeleteNotificationsByTypeActor(ctx context.Context, recipientID ulid.ULID, notifType entity.NotificationType, actorID ulid.ULID) (int64, error) {
	return r.q.DeleteNotificationsByTypeActor(ctx, sqlc.DeleteNotificationsByTypeActorParams{
		RecipientUserID: recipientID,
		Type:            sqlc.NotificationType(notifType),
		ActorUserID:     actorID,
	})
}
