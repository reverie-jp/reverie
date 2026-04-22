package handler

import (
	"context"

	"connectrpc.com/connect"

	notificationv1 "reverie.jp/reverie/internal/gen/pb/notification/v1"
	"reverie.jp/reverie/internal/modules/notification/adapter"
)

func (h *Handler) ListNotifications(ctx context.Context, req *connect.Request[notificationv1.ListNotificationsRequest]) (*connect.Response[notificationv1.ListNotificationsResponse], error) {
	input, err := adapter.FromListNotificationsRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.listNotifications.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListNotificationsResponse(output), nil
}
