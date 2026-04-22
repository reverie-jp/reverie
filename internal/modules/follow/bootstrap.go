package follow

import (
	"reverie.jp/reverie/internal/gen/pb/follow/v1/followv1connect"
	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	"reverie.jp/reverie/internal/modules/follow/handler"
	"reverie.jp/reverie/internal/modules/follow/usecase"
	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

func InitModule(followGateway followgw.Gateway, userGateway usergw.Gateway, notificationGateway notificationgw.Gateway) followv1connect.FollowServiceHandler {
	followUser := usecase.NewFollowUser(followGateway, userGateway, notificationGateway)
	unfollowUser := usecase.NewUnfollowUser(followGateway, userGateway, notificationGateway)
	listFollowingUsers := usecase.NewListFollowingUsers(followGateway, userGateway)
	listUserFollowers := usecase.NewListUserFollowers(followGateway, userGateway)

	return handler.New(followUser, unfollowUser, listFollowingUsers, listUserFollowers)
}
