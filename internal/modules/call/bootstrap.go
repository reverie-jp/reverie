package call

import (
	"time"

	"reverie.jp/reverie/internal/gen/pb/call/v1/callv1connect"
	"reverie.jp/reverie/internal/modules/call/handler"
	"reverie.jp/reverie/internal/modules/call/usecase"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/livekit"
)

func InitModule(lk *livekit.Client, userGateway usergw.Gateway, tokenTTL time.Duration) callv1connect.CallServiceHandler {
	joinCall := usecase.NewJoinCall(lk, userGateway, tokenTTL)
	return handler.New(joinCall)
}
