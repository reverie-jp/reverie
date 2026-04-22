package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func (r *RepositoryImpl) CreateNotification(ctx context.Context, params CreateNotificationParams) (*entity.Notification, error) {
	row, err := r.q.CreateNotification(ctx, sqlc.CreateNotificationParams{
		ID:              params.ID,
		RecipientUserID: params.RecipientUserID,
		Type:            sqlc.NotificationType(params.Type),
		ActorUserID:     params.ActorUserID,
		ResourceName:    params.ResourceName,
	})
	if err != nil {
		return nil, err
	}
	return mapper.ToNotification(&row), nil
}
