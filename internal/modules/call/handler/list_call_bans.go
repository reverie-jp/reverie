package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) ListCallBans(ctx context.Context, req *connect.Request[callv1.ListCallBansRequest]) (*connect.Response[callv1.ListCallBansResponse], error) {
	input, err := adapter.FromListCallBansRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.listCallBans.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListCallBansResponse(output), nil
}
