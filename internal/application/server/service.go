package server

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	"reverie.jp/reverie/internal/application/server/interceptor"
	"reverie.jp/reverie/internal/application/transaction"
	"reverie.jp/reverie/internal/config"
	accountv1 "reverie.jp/reverie/internal/gen/pb/account/v1"
	"reverie.jp/reverie/internal/gen/pb/account/v1/accountv1connect"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/gen/pb/user/v1/userv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/account"
	"reverie.jp/reverie/internal/modules/user"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/google"
	"reverie.jp/reverie/internal/platform/jwt"
)

type Service struct {
	Name                   string
	RegisterConnectHandler func(mux *http.ServeMux)
	RegisterGatewayHandler func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error
}

func initServices(cfg *config.Config, db *pgxpool.Pool, jwtManager *jwt.Manager) []Service {
	q := sqlc.New(db)
	tx := transaction.NewRunner(db)
	googleAuth := google.NewAuthClient(cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.RedirectURL)
	errorInterceptor := interceptor.ErrorInterceptor(cfg.Env)
	authInterceptor := interceptor.AuthInterceptor(jwtManager)
	userGateway := usergw.New(q)
	accountService := account.InitModule(q, userGateway, tx, googleAuth, jwtManager)
	userService := user.InitModule(userGateway)

	return []Service{
		{
			Name: accountv1connect.AccountServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(accountv1connect.NewAccountServiceHandler(
					accountService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
			RegisterGatewayHandler: func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
				return accountv1.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, addr, opts)
			},
		},
		{
			Name: userv1connect.UserServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(userv1connect.NewUserServiceHandler(
					userService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
			RegisterGatewayHandler: func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
				return userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addr, opts)
			},
		},
	}
}
