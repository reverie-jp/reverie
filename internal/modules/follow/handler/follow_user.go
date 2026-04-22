package handler

import (
	"context"

	"connectrpc.com/connect"

	followv1 "reverie.jp/reverie/internal/gen/pb/follow/v1"
	"reverie.jp/reverie/internal/modules/follow/adapter"
)

func (h *Handler) FollowUser(ctx context.Context, req *connect.Request[followv1.FollowUserRequest]) (*connect.Response[followv1.FollowUserResponse], error) {
	input, err := adapter.FromFollowUserRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.followUser.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToFollowUserResponse(output), nil
}
