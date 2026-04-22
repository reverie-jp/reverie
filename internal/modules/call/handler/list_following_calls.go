package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) ListFollowingCalls(ctx context.Context, req *connect.Request[callv1.ListFollowingCallsRequest]) (*connect.Response[callv1.ListFollowingCallsResponse], error) {
	input, err := adapter.FromListFollowingCallsRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.listFollowingCalls.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListFollowingCallsResponse(output), nil
}
