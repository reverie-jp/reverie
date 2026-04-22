package handler

import (
	"context"

	"connectrpc.com/connect"

	notificationv1 "reverie.jp/reverie/internal/gen/pb/notification/v1"
	"reverie.jp/reverie/internal/modules/notification/adapter"
)

func (h *Handler) MarkNotificationsRead(ctx context.Context, req *connect.Request[notificationv1.MarkNotificationsReadRequest]) (*connect.Response[notificationv1.MarkNotificationsReadResponse], error) {
	input, err := adapter.FromMarkNotificationsReadRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.markNotificationsRead.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToMarkNotificationsReadResponse(output), nil
}
