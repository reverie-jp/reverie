package handler

import (
	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/modules/account/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type Handler struct {
	accountv1connect.UnimplementedAccountServiceHandler
	socialLogin   *usecase.SocialLogin
	refreshToken  *usecase.RefreshToken
	deleteAccount *usecase.DeleteAccount
}

func New(
	socialLogin *usecase.SocialLogin,
	refreshToken *usecase.RefreshToken,
	deleteAccount *usecase.DeleteAccount,
) *Handler {
	return &Handler{
		socialLogin:   socialLogin,
		refreshToken:  refreshToken,
		deleteAccount: deleteAccount,
	}
}

func toProviderString(p accountv1.AuthProvider) (string, error) {
	switch p {
	case accountv1.AuthProvider_AUTH_PROVIDER_GOOGLE:
		return "google", nil
	default:
		return "", xerrors.ErrInvalidArgument.WithMessage("unsupported auth provider")
	}
}
