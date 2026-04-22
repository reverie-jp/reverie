package usecase

import (
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
)

type GetUserParticipatingCallOutput struct {
	// Nil if the user is not participating in any visible call.
	View *callgw.CallView
}
