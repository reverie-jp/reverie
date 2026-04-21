package adapter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"reverie.jp/reverie/internal/domain/entity"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	"reverie.jp/reverie/internal/platform/resourcename"
)

func ToCall(c *entity.Call, host *usergw.UserView) *callv1.Call {
	if c == nil {
		return nil
	}
	out := &callv1.Call{
		Name:       resourcename.FormatCall(c.ID),
		Host:       useradapter.ToUser(host),
		Visibility: toProtoVisibility(c.Visibility),
		CreateTime: timestamppb.New(c.CreateTime),
	}
	if c.EndTime != nil {
		out.EndTime = timestamppb.New(*c.EndTime)
	}
	return out
}
