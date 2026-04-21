package usecase

import (
	"reverie.jp/reverie/internal/domain/entity"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

type GetCallOutput struct {
	Call         *entity.Call
	Host         *usergw.UserView
	Participants []*CallParticipantView
}

type CallParticipantView struct {
	Participant          *entity.CallParticipant
	UserView             *usergw.UserView
	IsCurrentlyConnected bool
}
