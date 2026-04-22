package usecase

import (
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

type UnfollowUserOutput struct {
	View *usergw.UserView
}
