package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) KickCallParticipant(ctx context.Context, req *connect.Request[callv1.KickCallParticipantRequest]) (*connect.Response[callv1.KickCallParticipantResponse], error) {
	input, err := adapter.FromKickCallParticipantRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.kickCallParticipant.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToKickCallParticipantResponse(), nil
}
