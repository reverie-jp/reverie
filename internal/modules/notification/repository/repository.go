package repository

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type CreateNotificationParams struct {
	ID              ulid.ULID
	RecipientUserID ulid.ULID
	Type            entity.NotificationType
	ActorUserID     *ulid.ULID
	ResourceName    string
}

type Repository interface {
	CreateNotification(ctx context.Context, params CreateNotificationParams) (*entity.Notification, error)
	// CreateFanOutNotifications inserts many notifications that share actor /
	// type / resource_name in a single statement. Returns only newly-inserted
	// rows (dedup-conflicted ones are omitted).
	CreateFanOutNotifications(ctx context.Context, params CreateFanOutNotificationsParams) ([]*entity.Notification, error)
	ListNotificationsByRecipient(ctx context.Context, recipientID ulid.ULID, cursorID string, pageSize int32) ([]*entity.Notification, error)
	MarkNotificationsRead(ctx context.Context, recipientID ulid.ULID, ids []ulid.ULID) (int64, error)
	MarkAllNotificationsRead(ctx context.Context, recipientID ulid.ULID) (int64, error)
	CountUnreadNotifications(ctx context.Context, recipientID ulid.ULID) (int32, error)
}

type RepositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &RepositoryImpl{q: q}
}
