package user

import (
	"reverie.jp/reverie/internal/gen/pb/user/v1/userv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/user/handler"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/modules/user/usecase"
)

func InitModule(q *sqlc.Queries, userGateway usergw.Gateway) userv1connect.UserServiceHandler {
	getUser := usecase.NewGetUser(userGateway)
	updateUser := usecase.NewUpdateUser(userGateway)
	getUserSettings := usecase.NewGetUserSettings(userGateway)
	updateUserSettings := usecase.NewUpdateUserSettings(userGateway)
	followUser := usecase.NewFollowUser(userGateway)
	unfollowUser := usecase.NewUnfollowUser(userGateway)
	searchUsers := usecase.NewSearchUsers(userGateway)

	return handler.New(getUser, updateUser, getUserSettings, updateUserSettings, followUser, unfollowUser, searchUsers)
}
