package usecase

import (
	"reverie.jp/reverie/internal/domain/entity"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

type CallBanView struct {
	Ban  *entity.CallBan
	User *usergw.UserView
}

type ListCallBansOutput struct {
	Bans          []*CallBanView
	NextPageToken string
}
