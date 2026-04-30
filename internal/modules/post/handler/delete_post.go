package handler

import (
	"context"

	"connectrpc.com/connect"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/adapter"
)

func (h *Handler) DeletePost(ctx context.Context, req *connect.Request[postv1.DeletePostRequest]) (*connect.Response[postv1.DeletePostResponse], error) {
	input, userID, err := adapter.FromDeletePostRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := h.deletePost.Execute(ctx, input, userID); err != nil {
		return nil, err
	}
	return connect.NewResponse(&postv1.DeletePostResponse{}), nil
}
