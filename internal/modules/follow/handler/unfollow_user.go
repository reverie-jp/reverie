package handler

import (
	"context"

	"connectrpc.com/connect"

	followv1 "reverie.jp/reverie/internal/gen/pb/follow/v1"
	"reverie.jp/reverie/internal/modules/follow/adapter"
)

func (h *Handler) UnfollowUser(ctx context.Context, req *connect.Request[followv1.UnfollowUserRequest]) (*connect.Response[followv1.UnfollowUserResponse], error) {
	input, err := adapter.FromUnfollowUserRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.unfollowUser.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToUnfollowUserResponse(output), nil
}
