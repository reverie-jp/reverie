package usecase

import (
	"reverie.jp/reverie/internal/domain/entity"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

type TransferCallHostOutput struct {
	Call *entity.Call
	Host *usergw.UserView
}
