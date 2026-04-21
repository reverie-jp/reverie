package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) MuteCallParticipant(ctx context.Context, req *connect.Request[callv1.MuteCallParticipantRequest]) (*connect.Response[callv1.MuteCallParticipantResponse], error) {
	input, err := adapter.FromMuteCallParticipantRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.muteCallParticipant.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToMuteCallParticipantResponse(), nil
}
