package server

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"reverie.jp/reverie/internal/config"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/gen/pb/user/v1/userv1connect"
	"reverie.jp/reverie/internal/platform/jwt"
)

type Service struct {
	Name                   string
	RegisterConnectHandler func(mux *http.ServeMux)
	RegisterGatewayHandler func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error
}

func initServices(cfg *config.Config, db *pgxpool.Pool, jwtManager *jwt.Manager) []Service {
	return []Service{
		{
			Name: userv1connect.UserServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(userv1connect.NewUserServiceHandler(
					nil,                        // TODO: pass actual implementation
					connect.WithInterceptors(), // TODO: pass actual interceptors
				))
			},
			RegisterGatewayHandler: func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
				return userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addr, opts)
			},
		},
	}
}
