package usecase

import (
	"context"

	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListNotifications struct {
	notificationGateway notificationgw.Gateway
}

func NewListNotifications(notificationGateway notificationgw.Gateway) *ListNotifications {
	return &ListNotifications{notificationGateway: notificationGateway}
}

func (uc *ListNotifications) Execute(ctx context.Context, input ListNotificationsInput) (*ListNotificationsOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultListPageSize
	}
	if pageSize > maxListPageSize {
		pageSize = maxListPageSize
	}

	views, err := uc.notificationGateway.ListByRecipient(ctx, input.RequesterID, input.PageToken, pageSize)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}

	nextPageToken := ""
	if int32(len(views)) == pageSize && len(views) > 0 {
		last := views[len(views)-1]
		if last != nil && last.Notification != nil {
			nextPageToken = last.Notification.ID.String()
		}
	}

	return &ListNotificationsOutput{Views: views, NextPageToken: nextPageToken}, nil
}
