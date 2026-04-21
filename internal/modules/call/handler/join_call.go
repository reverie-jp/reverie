package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) JoinCall(ctx context.Context, req *connect.Request[callv1.JoinCallRequest]) (*connect.Response[callv1.JoinCallResponse], error) {
	input, err := adapter.FromJoinCallRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.joinCall.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToJoinCallResponse(output), nil
}
