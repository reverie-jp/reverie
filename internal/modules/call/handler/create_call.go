package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) CreateCall(ctx context.Context, req *connect.Request[callv1.CreateCallRequest]) (*connect.Response[callv1.CreateCallResponse], error) {
	input, err := adapter.FromCreateCallRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.createCall.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToCreateCallResponse(output), nil
}
