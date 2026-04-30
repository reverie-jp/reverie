package handler

import (
	"context"

	"connectrpc.com/connect"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/adapter"
)

func (h *Handler) LikePost(ctx context.Context, req *connect.Request[postv1.LikePostRequest]) (*connect.Response[postv1.LikePostResponse], error) {
	input, userID, err := adapter.FromLikePostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.likePost.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return adapter.ToLikePostResponse(output), nil
}

func (h *Handler) UnlikePost(ctx context.Context, req *connect.Request[postv1.UnlikePostRequest]) (*connect.Response[postv1.UnlikePostResponse], error) {
	input, userID, err := adapter.FromUnlikePostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.unlikePost.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return adapter.ToUnlikePostResponse(output), nil
}
