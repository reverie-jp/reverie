package adapter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/resourcename"
)

func ToUser(view *gateway.UserView) *userv1.User {
	if view == nil || view.User == nil {
		return nil
	}
	u := view.User
	return &userv1.User{
		Name:           resourcename.FormatUser(u.CustomID),
		CustomId:       u.CustomID,
		DisplayName:    u.DisplayName,
		Biography:      u.Biography,
		Location:       u.Location,
		Website:        u.Website,
		AvatarUrl:      u.AvatarURL,
		BannerUrl:      u.BannerURL,
		IsPrivate:      u.IsPrivate,
		OnlineStatus:   userv1.OnlineStatus_ONLINE_STATUS_UNSPECIFIED,
		FollowingCount: u.FollowingCount,
		FollowerCount:  u.FollowerCount,
		IsFollowing:    view.IsFollowing,
		IsFollowedBy:   view.IsFollowedBy,
		IsMe:           view.IsMe,
		CreateTime:     timestamppb.New(u.CreateTime),
	}
}
