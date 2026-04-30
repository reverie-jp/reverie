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
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/gen/pb/call/v1/callv1connect"
	"reverie.jp/reverie/internal/gen/pb/event/v1/eventv1connect"
	followv1 "reverie.jp/reverie/internal/gen/pb/follow/v1"
	"reverie.jp/reverie/internal/gen/pb/follow/v1/followv1connect"
	notificationv1 "reverie.jp/reverie/internal/gen/pb/notification/v1"
	"reverie.jp/reverie/internal/gen/pb/notification/v1/notificationv1connect"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/gen/pb/post/v1/postv1connect"
	presencev1 "reverie.jp/reverie/internal/gen/pb/presence/v1"
	"reverie.jp/reverie/internal/gen/pb/presence/v1/presencev1connect"
	timelinev1 "reverie.jp/reverie/internal/gen/pb/timeline/v1"
	"reverie.jp/reverie/internal/gen/pb/timeline/v1/timelinev1connect"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/gen/pb/user/v1/userv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/account"
	"reverie.jp/reverie/internal/modules/call"
	"reverie.jp/reverie/internal/modules/event"
	"reverie.jp/reverie/internal/modules/follow"
	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	"reverie.jp/reverie/internal/modules/notification"
	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	notificationrepo "reverie.jp/reverie/internal/modules/notification/repository"
	"reverie.jp/reverie/internal/modules/post"
	"reverie.jp/reverie/internal/modules/presence"
	"reverie.jp/reverie/internal/modules/timeline"
	"reverie.jp/reverie/internal/modules/user"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/events"
	"reverie.jp/reverie/internal/platform/google"
	"reverie.jp/reverie/internal/platform/jwt"
	"reverie.jp/reverie/internal/platform/livekit"
	"reverie.jp/reverie/internal/platform/ratelimit"
)

type Service struct {
	Name                   string
	RegisterConnectHandler func(mux *http.ServeMux)
	// RegisterGatewayHandler may be nil for streaming-only services that
	// cannot be proxied through grpc-gateway.
	RegisterGatewayHandler func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error
}

func initServices(cfg *config.Config, db *pgxpool.Pool, jwtManager *jwt.Manager, eventBus events.Bus, limiter ratelimit.Limiter) []Service {
	q := sqlc.New(db)
	tx := transaction.NewRunner(db)
	googleAuth := google.NewAuthClient(cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.RedirectURL)
	livekitClient := livekit.NewClient(cfg.LiveKit.URL, cfg.LiveKit.APIKey, cfg.LiveKit.APISecret)
	errorInterceptor := interceptor.ErrorInterceptor(cfg.Env)
	authInterceptor := interceptor.AuthInterceptor(jwtManager)
	followGateway := followgw.New(q)
	userGateway := usergw.New(q, followGateway)
	notificationGateway := notificationgw.New(notificationrepo.New(q), userGateway, eventBus, limiter)
	accountService := account.InitModule(q, userGateway, tx, googleAuth, jwtManager)
	userService := user.InitModule(userGateway)
	callService := call.InitModule(q, userGateway, followGateway, notificationGateway, livekitClient, cfg.LiveKit.TokenTTL)
	followService := follow.InitModule(followGateway, userGateway, notificationGateway)
	notificationService := notification.InitModule(notificationGateway)
	presenceService := presence.InitModule(userGateway)
	eventService := event.InitModule(eventBus)
	postService := post.InitModule(q, userGateway)
	timelineService := timeline.InitModule(q, userGateway)

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
		{
			Name: callv1connect.CallServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(callv1connect.NewCallServiceHandler(
					callService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
			RegisterGatewayHandler: func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
				return callv1.RegisterCallServiceHandlerFromEndpoint(ctx, mux, addr, opts)
			},
		},
		{
			Name: followv1connect.FollowServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(followv1connect.NewFollowServiceHandler(
					followService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
			RegisterGatewayHandler: func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
				return followv1.RegisterFollowServiceHandlerFromEndpoint(ctx, mux, addr, opts)
			},
		},
		{
			Name: notificationv1connect.NotificationServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(notificationv1connect.NewNotificationServiceHandler(
					notificationService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
			RegisterGatewayHandler: func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
				return notificationv1.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, addr, opts)
			},
		},
		{
			Name: presencev1connect.PresenceServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(presencev1connect.NewPresenceServiceHandler(
					presenceService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
			RegisterGatewayHandler: func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
				return presencev1.RegisterPresenceServiceHandlerFromEndpoint(ctx, mux, addr, opts)
			},
		},
		{
			// EventService is server-streaming only; grpc-gateway would not
			// bridge it to HTTP/JSON, so Connect alone serves it.
			Name: eventv1connect.EventServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(eventv1connect.NewEventServiceHandler(
					eventService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
		},
		{
			Name: postv1connect.PostServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(postv1connect.NewPostServiceHandler(
					postService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
			RegisterGatewayHandler: func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
				return postv1.RegisterPostServiceHandlerFromEndpoint(ctx, mux, addr, opts)
			},
		},
		{
			Name: timelinev1connect.TimelineServiceName,
			RegisterConnectHandler: func(mux *http.ServeMux) {
				mux.Handle(timelinev1connect.NewTimelineServiceHandler(
					timelineService,
					connect.WithInterceptors(errorInterceptor, authInterceptor),
				))
			},
			RegisterGatewayHandler: func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
				return timelinev1.RegisterTimelineServiceHandlerFromEndpoint(ctx, mux, addr, opts)
			},
		},
	}
}
