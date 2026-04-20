package user

import (
	"reverie.jp/reverie/internal/gen/pb/user/v1/userv1connect"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/modules/user/handler"
	"reverie.jp/reverie/internal/modules/user/usecase"
)

func InitModule(userGateway usergw.Gateway) userv1connect.UserServiceHandler {
	getUser := usecase.NewGetUser(userGateway)
	return handler.New(getUser)
}
