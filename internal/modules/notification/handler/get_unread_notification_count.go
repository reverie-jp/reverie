package handler

import (
	"context"

	"connectrpc.com/connect"

	notificationv1 "reverie.jp/reverie/internal/gen/pb/notification/v1"
	"reverie.jp/reverie/internal/modules/notification/adapter"
)

func (h *Handler) GetUnreadNotificationCount(ctx context.Context, req *connect.Request[notificationv1.GetUnreadNotificationCountRequest]) (*connect.Response[notificationv1.GetUnreadNotificationCountResponse], error) {
	input, err := adapter.FromGetUnreadNotificationCountRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.getUnreadNotificationCount.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToGetUnreadNotificationCountResponse(output), nil
}
