package usecase

import (
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
)

type ListCallBansOutput struct {
	Bans          []*callgw.CallBanView
	NextPageToken string
}
