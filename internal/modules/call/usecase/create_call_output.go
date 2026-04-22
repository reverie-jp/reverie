package usecase

import (
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
)

type CreateCallOutput struct {
	View *callgw.CallView
}
