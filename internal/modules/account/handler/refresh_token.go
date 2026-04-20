package handler

import (
	"context"

	"connectrpc.com/connect"

	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/modules/account/adapter"
)

func (h *Handler) RefreshToken(ctx context.Context, req *connect.Request[accountv1.RefreshTokenRequest]) (*connect.Response[accountv1.RefreshTokenResponse], error) {
	input, err := adapter.FromRefreshTokenRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	output, err := h.refreshToken.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return adapter.ToRefreshTokenResponse(output), nil
}
