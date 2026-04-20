package adapter

import (
	"context"

	"connectrpc.com/connect"

	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/modules/account/usecase"
)

func FromSocialLoginRequest(_ context.Context, req *connect.Request[accountv1.SocialLoginRequest]) (usecase.SocialLoginInput, error) {
	provider, err := toProviderString(req.Msg.Provider)
	if err != nil {
		return usecase.SocialLoginInput{}, err
	}
	return usecase.SocialLoginInput{
		Provider: provider,
		Code:     req.Msg.Code,
	}, nil
}

func ToSocialLoginResponse(output *usecase.SocialLoginOutput) *connect.Response[accountv1.SocialLoginResponse] {
	return connect.NewResponse(&accountv1.SocialLoginResponse{
		TokenPair:    toTokenPair(output.AccessToken, output.RefreshToken),
		IsNewAccount: output.IsNewAccount,
	})
}
