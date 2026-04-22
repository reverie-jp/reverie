package handler

import (
	"reverie.jp/reverie/internal/gen/pb/follow/v1/followv1connect"
	"reverie.jp/reverie/internal/modules/follow/usecase"
)

type Handler struct {
	followv1connect.UnimplementedFollowServiceHandler
	followUser         *usecase.FollowUser
	unfollowUser       *usecase.UnfollowUser
	listFollowingUsers *usecase.ListFollowingUsers
	listUserFollowers  *usecase.ListUserFollowers
}

func New(
	followUser *usecase.FollowUser,
	unfollowUser *usecase.UnfollowUser,
	listFollowingUsers *usecase.ListFollowingUsers,
	listUserFollowers *usecase.ListUserFollowers,
) *Handler {
	return &Handler{
		followUser:         followUser,
		unfollowUser:       unfollowUser,
		listFollowingUsers: listFollowingUsers,
		listUserFollowers:  listUserFollowers,
	}
}
