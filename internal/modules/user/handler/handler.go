package handler

import (
	"reverie.jp/reverie/internal/gen/pb/user/v1/userv1connect"
	"reverie.jp/reverie/internal/modules/user/usecase"
)

type Handler struct {
	userv1connect.UnimplementedUserServiceHandler
	getUser            *usecase.GetUser
	updateUser         *usecase.UpdateUser
	getUserSettings    *usecase.GetUserSettings
	updateUserSettings *usecase.UpdateUserSettings
	followUser         *usecase.FollowUser
	unfollowUser       *usecase.UnfollowUser
	searchUsers        *usecase.SearchUsers
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func New(
	getUser *usecase.GetUser,
	updateUser *usecase.UpdateUser,
	getUserSettings *usecase.GetUserSettings,
	updateUserSettings *usecase.UpdateUserSettings,
	followUser *usecase.FollowUser,
	unfollowUser *usecase.UnfollowUser,
	searchUsers *usecase.SearchUsers,
) *Handler {
	return &Handler{
		getUser:            getUser,
		updateUser:         updateUser,
		getUserSettings:    getUserSettings,
		updateUserSettings: updateUserSettings,
		followUser:         followUser,
		unfollowUser:       unfollowUser,
		searchUsers:        searchUsers,
	}
}
