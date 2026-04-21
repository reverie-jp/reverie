package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	callv1 "reverie.jp/reverie/internal/gen/pb/call/v1"
	"reverie.jp/reverie/internal/modules/call/usecase"
	"reverie.jp/reverie/internal/platform/resourcename"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromMuteCallParticipantRequest(ctx context.Context, req *connect.Request[callv1.MuteCallParticipantRequest]) (usecase.MuteCallParticipantInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.MuteCallParticipantInput{}, xerrors.ErrUnauthenticated
	}
	callID, identity, err := resourcename.ParseCallParticipant(req.Msg.Name)
	if err != nil {
		return usecase.MuteCallParticipantInput{}, err
	}
	return usecase.MuteCallParticipantInput{
		RequesterID: userID,
		CallID:      callID,
		Identity:    identity,
	}, nil
}

func ToMuteCallParticipantResponse() *connect.Response[callv1.MuteCallParticipantResponse] {
	return connect.NewResponse(&callv1.MuteCallParticipantResponse{})
}
