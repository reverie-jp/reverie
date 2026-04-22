package usecase

import (
	"context"

	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type MarkNotificationsRead struct {
	notificationGateway notificationgw.Gateway
}

func NewMarkNotificationsRead(notificationGateway notificationgw.Gateway) *MarkNotificationsRead {
	return &MarkNotificationsRead{notificationGateway: notificationGateway}
}

func (uc *MarkNotificationsRead) Execute(ctx context.Context, input MarkNotificationsReadInput) (*MarkNotificationsReadOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	var (
		count int32
		err   error
	)
	if len(input.IDs) == 0 {
		count, err = uc.notificationGateway.MarkAllRead(ctx, input.RequesterID)
	} else {
		count, err = uc.notificationGateway.MarkRead(ctx, input.RequesterID, input.IDs)
	}
	if err != nil {
		return nil, xerrors.ErrInternal.WithCause(err)
	}
	return &MarkNotificationsReadOutput{MarkedCount: count}, nil
}
