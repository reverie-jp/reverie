package handler

import (
	"context"

	"connectrpc.com/connect"

	followv1 "reverie.jp/reverie/internal/gen/pb/follow/v1"
	"reverie.jp/reverie/internal/modules/follow/adapter"
)

func (h *Handler) ListFollowingUsers(ctx context.Context, req *connect.Request[followv1.ListFollowingUsersRequest]) (*connect.Response[followv1.ListFollowingUsersResponse], error) {
	input, err := adapter.FromListFollowingUsersRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.listFollowingUsers.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToListFollowingUsersResponse(output), nil
}
