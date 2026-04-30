package handler

import (
	"context"

	"connectrpc.com/connect"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/adapter"
)

func (h *Handler) ListPostLikes(ctx context.Context, req *connect.Request[postv1.ListPostLikesRequest]) (*connect.Response[postv1.ListPostLikesResponse], error) {
	params, err := adapter.FromListPostLikesRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	views, err := h.listPostLikes.Execute(ctx, params.AuthorCustomID, params.ShortID, params.UserID, params.PageSize)
	if err != nil {
		return nil, err
	}
	return adapter.ToListPostLikesResponse(views), nil
}
