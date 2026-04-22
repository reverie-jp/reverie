package adapter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	callgw "reverie.jp/reverie/internal/modules/call/gateway"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	"reverie.jp/reverie/internal/platform/resourcename"
)

func ToCall(view *callgw.CallView) *callv1.Call {
	if view == nil || view.Call == nil {
		return nil
	}
	c := view.Call
	out := &callv1.Call{
		Name:       resourcename.FormatCall(c.ID),
		Host:       useradapter.ToUser(view.Host),
		Visibility: toProtoVisibility(c.Visibility),
		Title:      c.Title,
		CreateTime: timestamppb.New(c.CreateTime),
	}
	if c.EndTime != nil {
		out.EndTime = timestamppb.New(*c.EndTime)
	}
	if len(view.ActiveParticipants) > 0 {
		out.ActiveParticipants = make([]*callv1.CallParticipant, 0, len(view.ActiveParticipants))
		for _, pv := range view.ActiveParticipants {
			if p := ToCallParticipant(pv); p != nil {
				out.ActiveParticipants = append(out.ActiveParticipants, p)
			}
		}
	}
	return out
}
