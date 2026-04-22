package mapper

import (
	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
)

func ToNotification(row *sqlc.Notification) *entity.Notification {
	if row == nil {
		return nil
	}
	return &entity.Notification{
		ID:              row.ID,
		RecipientUserID: row.RecipientUserID,
		Type:            entity.NotificationType(row.Type),
		ActorUserID:     row.ActorUserID,
		ResourceName:    row.ResourceName,
		ReadTime:        row.ReadTime,
		CreateTime:      row.CreateTime,
	}
}
