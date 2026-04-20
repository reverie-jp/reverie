package handler

import (
	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/modules/account/usecase"
)

type Handler struct {
	accountv1connect.UnimplementedAccountServiceHandler
	socialLogin   *usecase.SocialLogin
	refreshToken  *usecase.RefreshToken
	logout        *usecase.Logout
	deleteAccount *usecase.DeleteAccount
}

func New(
	socialLogin *usecase.SocialLogin,
	refreshToken *usecase.RefreshToken,
	logout *usecase.Logout,
	deleteAccount *usecase.DeleteAccount,
) *Handler {
	return &Handler{
		socialLogin:   socialLogin,
		refreshToken:  refreshToken,
		logout:        logout,
		deleteAccount: deleteAccount,
	}
}
