package handler

import (
	"context"

	"connectrpc.com/connect"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/adapter"
)

func (h *Handler) TransferCallHost(ctx context.Context, req *connect.Request[callv1.TransferCallHostRequest]) (*connect.Response[callv1.TransferCallHostResponse], error) {
	input, err := adapter.FromTransferCallHostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.transferCallHost.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToTransferCallHostResponse(output), nil
}
