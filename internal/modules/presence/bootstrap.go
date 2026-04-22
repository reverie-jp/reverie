package presence

import (
	"reverie.jp/reverie/internal/gen/pb/presence/v1/presencev1connect"
	"reverie.jp/reverie/internal/modules/presence/handler"
	"reverie.jp/reverie/internal/modules/presence/usecase"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

func InitModule(userGateway usergw.Gateway) presencev1connect.PresenceServiceHandler {
	heartbeat := usecase.NewHeartbeat(userGateway)
	return handler.New(heartbeat)
}
