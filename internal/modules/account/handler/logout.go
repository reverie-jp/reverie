package handler

import (
	"context"

	"connectrpc.com/connect"

	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
)

func (h *Handler) Logout(ctx context.Context, req *connect.Request[accountv1.LogoutRequest]) (*connect.Response[accountv1.LogoutResponse], error) {
	// With stateless JWTs, logout is handled client-side by discarding tokens.
	// Future: add token blacklist for immediate revocation.
	return connect.NewResponse(&accountv1.LogoutResponse{}), nil
}
