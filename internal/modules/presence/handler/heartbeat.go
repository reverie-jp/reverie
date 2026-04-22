package handler

import (
	"context"

	"connectrpc.com/connect"

	presencev1 "reverie.jp/reverie/internal/gen/pb/presence/v1"
	"reverie.jp/reverie/internal/modules/presence/adapter"
)

func (h *Handler) Heartbeat(ctx context.Context, req *connect.Request[presencev1.HeartbeatRequest]) (*connect.Response[presencev1.HeartbeatResponse], error) {
	input, err := adapter.FromHeartbeatRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.heartbeat.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToHeartbeatResponse(output), nil
}
