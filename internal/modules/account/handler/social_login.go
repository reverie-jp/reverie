package handler

import (
	"context"

	"connectrpc.com/connect"

	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/modules/account/usecase"
)

func (h *Handler) SocialLogin(ctx context.Context, req *connect.Request[accountv1.SocialLoginRequest]) (*connect.Response[accountv1.SocialLoginResponse], error) {
	provider, err := toProviderString(req.Msg.Provider)
	if err != nil {
		return nil, err
	}

	output, err := h.socialLogin.Execute(ctx, usecase.SocialLoginInput{
		Provider: provider,
		Code:     req.Msg.Code,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.SocialLoginResponse{
		TokenPair: &accountv1.TokenPair{
			AccessToken:  output.AccessToken,
			RefreshToken: output.RefreshToken,
		},
		IsNewAccount: output.IsNewAccount,
	}), nil
}
