package usecase

import (
	"reverie.jp/reverie/internal/domain/entity"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

type GetUserParticipatingCallOutput struct {
	// Nil if the user is not participating in any visible call.
	Call *entity.Call
	Host *usergw.UserView
}
