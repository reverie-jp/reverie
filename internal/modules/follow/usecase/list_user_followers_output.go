package usecase

import (
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

type ListUserFollowersOutput struct {
	Views         []*usergw.UserView
	NextPageToken string
}
