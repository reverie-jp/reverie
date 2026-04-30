package handler

import (
	"context"

	"connectrpc.com/connect"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/adapter"
)

func (h *Handler) ListPostReplies(ctx context.Context, req *connect.Request[postv1.ListPostRepliesRequest]) (*connect.Response[postv1.ListPostRepliesResponse], error) {
	input, userID, err := adapter.FromListPostRepliesRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	outputs, err := h.listPostReplies.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return adapter.ToListPostRepliesResponse(outputs), nil
}
