package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) EndCall(ctx context.Context, req *connect.Request[callv1.EndCallRequest]) (*connect.Response[callv1.EndCallResponse], error) {
	input, err := adapter.FromEndCallRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.endCall.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToEndCallResponse(), nil
}
