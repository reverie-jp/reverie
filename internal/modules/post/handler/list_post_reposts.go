package handler

import (
	"context"

	"connectrpc.com/connect"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/adapter"
)

func (h *Handler) ListPostReposts(ctx context.Context, req *connect.Request[postv1.ListPostRepostsRequest]) (*connect.Response[postv1.ListPostRepostsResponse], error) {
	input, userID, err := adapter.FromListPostRepostsRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	outputs, err := h.listPostReposts.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return adapter.ToListPostRepostsResponse(outputs), nil
}
