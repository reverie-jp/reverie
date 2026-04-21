package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) UpdateCall(ctx context.Context, req *connect.Request[callv1.UpdateCallRequest]) (*connect.Response[callv1.UpdateCallResponse], error) {
	input, err := adapter.FromUpdateCallRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.updateCall.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToUpdateCallResponse(output), nil
}
