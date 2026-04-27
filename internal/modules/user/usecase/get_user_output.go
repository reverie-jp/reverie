package usecase

import usergw "reverie.jp/reverie/internal/modules/user/gateway"

type GetUserOutput struct {
	View *usergw.UserView
}
