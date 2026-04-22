package usecase

import (
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
)

type GetCallOutput struct {
	View         *callgw.CallView
	Participants []*callgw.CallParticipantView
}
