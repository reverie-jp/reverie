package handler

import (
	"context"

	"connectrpc.com/connect"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/adapter"
)

func (h *Handler) ListUserLikedPosts(ctx context.Context, req *connect.Request[postv1.ListUserLikedPostsRequest]) (*connect.Response[postv1.ListUserLikedPostsResponse], error) {
	input, userID, err := adapter.FromListUserLikedPostsRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	outputs, err := h.listUserLikedPosts.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return adapter.ToListUserLikedPostsResponse(outputs), nil
}
