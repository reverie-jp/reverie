package usecase

import (
	"context"

	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetUnreadNotificationCount struct {
	notificationGateway notificationgw.Gateway
}

func NewGetUnreadNotificationCount(notificationGateway notificationgw.Gateway) *GetUnreadNotificationCount {
	return &GetUnreadNotificationCount{notificationGateway: notificationGateway}
}

func (uc *GetUnreadNotificationCount) Execute(ctx context.Context, input GetUnreadNotificationCountInput) (*GetUnreadNotificationCountOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	count, err := uc.notificationGateway.CountUnread(ctx, input.RequesterID)
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	return &GetUnreadNotificationCountOutput{Count: count}, nil
}
