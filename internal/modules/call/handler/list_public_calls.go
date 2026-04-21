package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) ListPublicCalls(ctx context.Context, req *connect.Request[callv1.ListPublicCallsRequest]) (*connect.Response[callv1.ListPublicCallsResponse], error) {
	input := adapter.FromListPublicCallsRequest(ctx, req)
	output, err := h.listPublicCalls.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListPublicCallsResponse(output), nil
}
