package adapter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	"reverie.jp/reverie/internal/platform/resourcename"
)

func ToCallBan(view *usecase.CallBanView) *callv1.CallBan {
	if view == nil || view.Ban == nil {
		return nil
	}
	b := view.Ban
	return &callv1.CallBan{
		Name:       resourcename.FormatCallBan(b.CallID, b.UserID),
		User:       useradapter.ToUser(view.User),
		CreateTime: timestamppb.New(b.CreateTime),
	}
}
