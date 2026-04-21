package usecase

import (
	"reverie.jp/reverie/internal/domain/entity"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type ListPublicCallsOutput struct {
	Calls         []*entity.Call
	HostsByID     map[ulid.ULID]*usergw.UserView
	NextPageToken string
}
