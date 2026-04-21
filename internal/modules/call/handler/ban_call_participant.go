package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) BanCallParticipant(ctx context.Context, req *connect.Request[callv1.BanCallParticipantRequest]) (*connect.Response[callv1.BanCallParticipantResponse], error) {
	input, err := adapter.FromBanCallParticipantRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.banCallParticipant.Execute(ctx, input); err != nil {
		return nil, err
	}
	return adapter.ToBanCallParticipantResponse(), nil
}
