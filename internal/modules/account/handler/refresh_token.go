package handler

import (
	"context"

	"connectrpc.com/connect"

	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/modules/account/usecase"
)

func (h *Handler) RefreshToken(ctx context.Context, req *connect.Request[accountv1.RefreshTokenRequest]) (*connect.Response[accountv1.RefreshTokenResponse], error) {
	output, err := h.refreshToken.Execute(ctx, usecase.RefreshTokenInput{
		RefreshToken: req.Msg.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.RefreshTokenResponse{
		TokenPair: &accountv1.TokenPair{
			AccessToken:  output.AccessToken,
			RefreshToken: output.RefreshToken,
		},
	}), nil
}
