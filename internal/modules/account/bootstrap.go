package account

import (
	"reverie.jp/reverie/internal/application/transaction"
	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/account/handler"
	"reverie.jp/reverie/internal/modules/account/repository"
	"reverie.jp/reverie/internal/modules/account/usecase"
	"reverie.jp/reverie/internal/platform/google"
	"reverie.jp/reverie/internal/platform/jwt"
)

func InitModule(q *sqlc.Queries, tx transaction.Runner, googleAuth *google.AuthClient, jwtManager *jwt.Manager) accountv1connect.AccountServiceHandler {
	repo := repository.NewRepository(q)

	socialLogin := usecase.NewSocialLogin(repo, tx, googleAuth, jwtManager)
	refreshToken := usecase.NewRefreshToken(jwtManager)
	getAccount := usecase.NewGetAccount(repo)
	deleteAccount := usecase.NewDeleteAccount(repo)

	return handler.NewHandler(socialLogin, refreshToken, getAccount, deleteAccount)
}
