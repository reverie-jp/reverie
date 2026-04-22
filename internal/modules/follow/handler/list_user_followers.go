package handler

import (
	"context"

	"connectrpc.com/connect"

	followv1 "reverie.jp/reverie/internal/gen/pb/follow/v1"
	"reverie.jp/reverie/internal/modules/follow/adapter"
)

func (h *Handler) ListUserFollowers(ctx context.Context, req *connect.Request[followv1.ListUserFollowersRequest]) (*connect.Response[followv1.ListUserFollowersResponse], error) {
	input, err := adapter.FromListUserFollowersRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.listUserFollowers.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListUserFollowersResponse(output), nil
}
