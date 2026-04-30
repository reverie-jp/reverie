package handler

import (
	"context"

	"connectrpc.com/connect"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/adapter"
)

func (h *Handler) CreatePost(ctx context.Context, req *connect.Request[postv1.CreatePostRequest]) (*connect.Response[postv1.CreatePostResponse], error) {
	input, userID, err := adapter.FromCreatePostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.createPost.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return adapter.ToCreatePostResponse(output), nil
}
