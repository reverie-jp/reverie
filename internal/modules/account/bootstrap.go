package account

import (
	"reverie.jp/reverie/internal/application/transaction"
	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/account/handler"
	"reverie.jp/reverie/internal/modules/account/usecase"
	"reverie.jp/reverie/internal/platform/google"
	"reverie.jp/reverie/internal/platform/jwt"
)

func InitModule(q *sqlc.Queries, tx transaction.Runner, googleAuth *google.AuthClient, jwtManager *jwt.Manager) accountv1connect.AccountServiceHandler {
	socialLogin := usecase.NewSocialLogin(q, tx, googleAuth, jwtManager)
	refreshToken := usecase.NewRefreshToken(jwtManager)
	getAccount := usecase.NewGetAccount(q)
	deleteAccount := usecase.NewDeleteAccount(q)

	return handler.NewAccountHandler(socialLogin, refreshToken, getAccount, deleteAccount)
}
