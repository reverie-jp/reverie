package usecase

import (
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

type ListFollowingUsersOutput struct {
	Views         []*usergw.UserView
	NextPageToken string
}
