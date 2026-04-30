package handler

import (
	"context"

	"connectrpc.com/connect"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/adapter"
)

func (h *Handler) GetPost(ctx context.Context, req *connect.Request[postv1.GetPostRequest]) (*connect.Response[postv1.GetPostResponse], error) {
	input, userID, err := adapter.FromGetPostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.getPost.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return adapter.ToGetPostResponse(output), nil
}
