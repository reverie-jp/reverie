package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) UnmuteCallParticipant(ctx context.Context, req *connect.Request[callv1.UnmuteCallParticipantRequest]) (*connect.Response[callv1.UnmuteCallParticipantResponse], error) {
	input, err := adapter.FromUnmuteCallParticipantRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.unmuteCallParticipant.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToUnmuteCallParticipantResponse(), nil
}
