package usecase

import (
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
)

type ListPublicCallsOutput struct {
	Views         []*callgw.CallView
	NextPageToken string
}
