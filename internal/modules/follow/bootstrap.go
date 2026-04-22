package follow

import (
	"reverie.jp/reverie/internal/gen/pb/follow/v1/followv1connect"
	followgw "reverie.jp/reverie/internal/modules/follow/gateway"
	"reverie.jp/reverie/internal/modules/follow/handler"
	"reverie.jp/reverie/internal/modules/follow/usecase"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

func InitModule(followGateway followgw.Gateway, userGateway usergw.Gateway) followv1connect.FollowServiceHandler {
	followUser := usecase.NewFollowUser(followGateway, userGateway)
	unfollowUser := usecase.NewUnfollowUser(followGateway, userGateway)
	listFollowingUsers := usecase.NewListFollowingUsers(followGateway, userGateway)
	listUserFollowers := usecase.NewListUserFollowers(followGateway, userGateway)

	return handler.New(followUser, unfollowUser, listFollowingUsers, listUserFollowers)
}
