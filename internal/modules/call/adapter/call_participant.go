package adapter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	"reverie.jp/reverie/internal/platform/resourcename"
)

func ToCallParticipant(view *usecase.CallParticipantView) *callv1.CallParticipant {
	if view == nil || view.Participant == nil {
		return nil
	}
	p := view.Participant
	return &callv1.CallParticipant{
		Name:                 resourcename.FormatCallParticipant(p.CallID, p.ParticipantIdentity),
		User:                 useradapter.ToUser(view.UserView),
		DisplayName:          p.DisplayName,
		FirstJoinTime:        timestamppb.New(p.FirstJoinTime),
		IsCurrentlyConnected: view.IsCurrentlyConnected,
	}
}
