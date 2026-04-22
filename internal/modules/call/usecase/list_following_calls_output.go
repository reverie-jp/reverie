package usecase

import (
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
)

type ListFollowingCallsOutput struct {
	Views         []*callgw.CallView
	NextPageToken string
}
