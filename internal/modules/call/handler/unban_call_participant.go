package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) UnbanCallParticipant(ctx context.Context, req *connect.Request[callv1.UnbanCallParticipantRequest]) (*connect.Response[callv1.UnbanCallParticipantResponse], error) {
	input, err := adapter.FromUnbanCallParticipantRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.unbanCallParticipant.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToUnbanCallParticipantResponse(), nil
}
