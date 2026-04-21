package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) LeaveCall(ctx context.Context, req *connect.Request[callv1.LeaveCallRequest]) (*connect.Response[callv1.LeaveCallResponse], error) {
	input, err := adapter.FromLeaveCallRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.leaveCall.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToLeaveCallResponse(), nil
}
