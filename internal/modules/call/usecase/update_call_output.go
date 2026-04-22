package usecase

import (
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
)

type UpdateCallOutput struct {
	View *callgw.CallView
}
