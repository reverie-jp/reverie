package account

import (
	"reverie.jp/reverie/internal/application/transaction"
	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/account/handler"
	accountrepo "reverie.jp/reverie/internal/modules/account/repository"
	"reverie.jp/reverie/internal/modules/account/usecase"
	userrepo "reverie.jp/reverie/internal/modules/user/repository"
	"reverie.jp/reverie/internal/platform/google"
	"reverie.jp/reverie/internal/platform/jwt"
)

func InitModule(q *sqlc.Queries, tx transaction.Runner, googleAuth *google.AuthClient, jwtManager *jwt.Manager) accountv1connect.AccountServiceHandler {
	accountRepo := accountrepo.NewRepository(q)
	userRepo := userrepo.NewRepository(q)

	socialLogin := usecase.NewSocialLogin(accountRepo, userRepo, tx, googleAuth, jwtManager)
	refreshToken := usecase.NewRefreshToken(jwtManager)
	getAccount := usecase.NewGetAccount(userRepo)
	deleteAccount := usecase.NewDeleteAccount(userRepo)

	return handler.NewHandler(socialLogin, refreshToken, getAccount, deleteAccount)
}
