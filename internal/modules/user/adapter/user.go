package adapter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/user/gateway"
)

func ToUser(view *gateway.UserView) *userv1.User {
	if view == nil || view.User == nil {
		return nil
	}
	u := view.User
	bio := u.Biography
	return &userv1.User{
		Id:             u.ID.String(),
		CustomId:       u.CustomID,
		DisplayName:    u.DisplayName,
		Biography:      &bio,
		IsPrivate:      u.IsPrivate,
		IsMe:           view.IsMe,
		IsFollowing:    view.IsFollowing,
		IsFollowedBy:   view.IsFollowedBy,
		FollowerCount:  int32(view.FollowerCount),
		FollowingCount: int32(view.FollowingCount),
		CreateTime:     timestamppb.New(u.CreateTime),
	}
}
