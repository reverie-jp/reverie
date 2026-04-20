package account

import (
	"reverie.jp/reverie/internal/application/transaction"
	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/account/handler"
	accountrepo "reverie.jp/reverie/internal/modules/account/repository"
	"reverie.jp/reverie/internal/modules/account/usecase"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/google"
	"reverie.jp/reverie/internal/platform/jwt"
)

func InitModule(q *sqlc.Queries, userGateway usergw.Gateway, tx transaction.Runner, googleAuth *google.AuthClient, jwtManager *jwt.Manager) accountv1connect.AccountServiceHandler {
	accountRepo := accountrepo.New(q)

	socialLogin := usecase.NewSocialLogin(accountRepo, userGateway, tx, googleAuth, jwtManager)
	refreshToken := usecase.NewRefreshToken(jwtManager)
	deleteAccount := usecase.NewDeleteAccount(userGateway)

	return handler.New(socialLogin, refreshToken, deleteAccount)
}
