package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) GetUserParticipatingCall(ctx context.Context, req *connect.Request[callv1.GetUserParticipatingCallRequest]) (*connect.Response[callv1.GetUserParticipatingCallResponse], error) {
	input, err := adapter.FromGetUserParticipatingCallRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.getUserParticipatingCall.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToGetUserParticipatingCallResponse(output), nil
}
