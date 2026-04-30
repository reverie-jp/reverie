package handler

import (
	"context"

	"connectrpc.com/connect"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/adapter"
)

func (h *Handler) ListUserPosts(ctx context.Context, req *connect.Request[postv1.ListUserPostsRequest]) (*connect.Response[postv1.ListUserPostsResponse], error) {
	input, userID, err := adapter.FromListUserPostsRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	outputs, err := h.listUserPosts.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return adapter.ToListUserPostsResponse(outputs), nil
}
