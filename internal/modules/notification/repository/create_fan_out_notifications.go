package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/domain/mapper"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type CreateFanOutNotificationsParams struct {
	IDs          []ulid.ULID
	RecipientIDs []ulid.ULID
	Type         entity.NotificationType
	ActorUserID  *ulid.ULID
	ResourceName string
}

func (r *RepositoryImpl) CreateFanOutNotifications(ctx context.Context, params CreateFanOutNotificationsParams) ([]*entity.Notification, error) {
	if len(params.IDs) == 0 {
		return []*entity.Notification{}, nil
	}
	idStrs := make([]string, len(params.IDs))
	recipStrs := make([]string, len(params.RecipientIDs))
	for i, id := range params.IDs {
		idStrs[i] = id.String()
	}
	for i, id := range params.RecipientIDs {
		recipStrs[i] = id.String()
	}
	rows, err := r.q.CreateFanOutNotifications(ctx, sqlc.CreateFanOutNotificationsParams{
		Type:         sqlc.NotificationType(params.Type),
		ActorUserID:  params.ActorUserID,
		ResourceName: params.ResourceName,
		Ids:          idStrs,
		RecipientIds: recipStrs,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*entity.Notification, 0, len(rows))
	for i := range rows {
		out = append(out, mapper.ToNotification(&rows[i]))
	}
	return out, nil
}
