package usecase

import (
	"reverie.jp/reverie/internal/modules/user/gateway"
)

type GetUserOutput struct {
	View *gateway.UserView
}
