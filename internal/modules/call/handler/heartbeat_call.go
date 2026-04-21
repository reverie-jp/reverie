package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) HeartbeatCall(ctx context.Context, req *connect.Request[callv1.HeartbeatCallRequest]) (*connect.Response[callv1.HeartbeatCallResponse], error) {
	input, err := adapter.FromHeartbeatCallRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.heartbeatCall.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToHeartbeatCallResponse(), nil
}
