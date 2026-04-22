package usecase

import (
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

type FollowUserOutput struct {
	View *usergw.UserView
}
