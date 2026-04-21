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

func FromBanCallParticipantRequest(ctx context.Context, req *connect.Request[callv1.BanCallParticipantRequest]) (usecase.BanCallParticipantInput, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.BanCallParticipantInput{}, xerrors.ErrUnauthenticated
	}
	callID, identity, err := resourcename.ParseCallParticipant(req.Msg.Name)
	if err != nil {
		return usecase.BanCallParticipantInput{}, err
	}
	return usecase.BanCallParticipantInput{
		RequesterID: userID,
		CallID:      callID,
		Identity:    identity,
	}, nil
}

func ToBanCallParticipantResponse() *connect.Response[callv1.BanCallParticipantResponse] {
	return connect.NewResponse(&callv1.BanCallParticipantResponse{})
}
