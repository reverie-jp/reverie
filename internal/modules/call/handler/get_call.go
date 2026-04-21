package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) GetCall(ctx context.Context, req *connect.Request[callv1.GetCallRequest]) (*connect.Response[callv1.GetCallResponse], error) {
	input, err := adapter.FromGetCallRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.getCall.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToGetCallResponse(output), nil
}
