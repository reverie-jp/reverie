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

func FromKickCallParticipantRequest(ctx context.Context, req *connect.Request[callv1.KickCallParticipantRequest]) (usecase.KickCallParticipantInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.KickCallParticipantInput{}, xerrors.ErrUnauthenticated
	}
	callID, identity, err := resourcename.ParseCallParticipant(req.Msg.Name)
	if err != nil {
		return usecase.KickCallParticipantInput{}, err
	}
	return usecase.KickCallParticipantInput{
		RequesterID: userID,
		CallID:      callID,
		Identity:    identity,
	}, nil
}

func ToKickCallParticipantResponse() *connect.Response[callv1.KickCallParticipantResponse] {
	return connect.NewResponse(&callv1.KickCallParticipantResponse{})
}
