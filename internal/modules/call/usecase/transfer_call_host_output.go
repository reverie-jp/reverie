package usecase

import (
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
)

type TransferCallHostOutput struct {
	View *callgw.CallView
}
